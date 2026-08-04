<?php

namespace App\Services\JobManagements;

use App\Models\JobManagement\JobManagementAsset;
use App\Models\JobManagement\JobManagementEducationExperiences;
use App\Models\JobManagement\JobManagementFinancial;
use App\Models\JobManagement\JobManagementPotencyCompetency;
use App\Models\JobManagement\JobManagementScore;
use App\Models\JobManagement\JobManagementRelationship;
use App\Models\JobManagement\JobManagementSubordinateControl;
use App\Models\JobManagement\JobManagementWorkingActivity;
use App\Models\JobManagement\JobManagementWorkingRisk;
use Illuminate\Support\Collection;

class JobValueCalculator
{
    private const MAP_DEFAULT = [
        1 => 1,
        2 => 3,
        3 => 6,
        4 => 10,
        5 => 15,
    ];

    private const MAP_EXTENDED = [
        1 => 1,
        2 => 3,
        3 => 6,
        4 => 10,
        5 => 15,
        6 => 21,
        7 => 28,
        8 => 36,
    ];

    private const MAP_LINEAR_5 = [
        1 => 1,
        2 => 2,
        3 => 3,
        4 => 4,
        5 => 5,
    ];

    private const MAP_LINEAR_8 = [
        1 => 1,
        2 => 2,
        3 => 3,
        4 => 4,
        5 => 5,
        6 => 6,
        7 => 7,
        8 => 8,
    ];

    private const MAP_COMMUNICATION = [
        1 => 1,
        2 => 3,
        3 => 6,
    ];

    /**
     * Calculate job value for a single organization.
     */
    public function calculate(int $organizationId): array
    {
        $results = $this->calculateForOrganizationIds([$organizationId]);

        return $results[$organizationId] ?? $this->emptyResult();
    }

    /**
     * Calculate job values for multiple organizations at once.
     *
     * @param  array<int, int|string>  $organizationIds
     * @return array<int, array>
     */
    public function calculateForOrganizationIds(array $organizationIds): array
    {
        $ids = collect($organizationIds)
            ->filter()
            ->map(fn($value) => (int) $value)
            ->unique()
            ->values();

        if ($ids->isEmpty()) {
            return [];
        }

        $educationExperiences = JobManagementEducationExperiences::with(['educationGrade', 'experienceGrade'])
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $potencyCompetencies = JobManagementPotencyCompetency::with('jobManagementValue')
            ->whereIn('organization_id', $ids)
            ->get()
            ->groupBy('organization_id');

        $financialRecords = JobManagementFinancial::with(['cashValue', 'authorityValue', 'impactValue'])
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $assetRecords = JobManagementAsset::with(['assetValue', 'assetAuthority'])
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $subordinateRecords = JobManagementSubordinateControl::with('value')
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $relationshipRecords = JobManagementRelationship::with(['scopeValue', 'frequencyValue'])
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $activityRecords = JobManagementWorkingActivity::with('value')
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $riskRecords = JobManagementWorkingRisk::with(['environmentValue', 'hazardValue'])
            ->whereIn('organization_id', $ids)
            ->get()
            ->keyBy('organization_id');

        $results = [];

        foreach ($ids as $organizationId) {
            $results[$organizationId] = $this->calculateForSingleOrganization(
                $educationExperiences->get($organizationId),
                $potencyCompetencies->get($organizationId, collect()),
                $financialRecords->get($organizationId),
                $assetRecords->get($organizationId),
                $subordinateRecords->get($organizationId),
                $relationshipRecords->get($organizationId),
                $activityRecords->get($organizationId),
                $riskRecords->get($organizationId)
            );
        }

        $this->persistResults($results);

        return $results;
    }

    private function calculateForSingleOrganization(
        ?JobManagementEducationExperiences $educationExperience,
        Collection $potencyCompetencies,
        ?JobManagementFinancial $financial,
        ?JobManagementAsset $asset,
        ?JobManagementSubordinateControl $subordinate,
        ?JobManagementRelationship $relationship,
        ?JobManagementWorkingActivity $activity,
        ?JobManagementWorkingRisk $risk
    ): array {
        $educationComponent = $this->calculateEducationExperience($educationExperience);
        $potentialComponent = $this->calculatePotentials($potencyCompetencies);
        $competencyComponent = $this->calculateCompetencies($potencyCompetencies);
        $problemSolvingComponent = $this->calculateProblemSolving($potencyCompetencies);
        $financialComponent = $this->calculateFinancialAuthority($financial);
        $assetComponent = $this->calculateAssetAuthority($asset);
        $subordinateComponent = $this->calculateSubordinateControl($subordinate);
        $relationshipComponent = $this->calculateWorkScope($relationship);
        $activityComponent = $this->calculateWorkActivity($activity);
        $riskComponent = $this->calculateWorkRisk($risk);

        $competencyBaseScore = $competencyComponent['score'] ?? 0;
        $potentialScore = $potentialComponent['score'] ?? 0;
        $problemSolvingScore = $problemSolvingComponent['score'] ?? 0;
        $competencyAggregateScore = $competencyBaseScore + $potentialScore + $problemSolvingScore;

        $competencyComponent['base_score'] = $competencyBaseScore;
        $competencyComponent['potential_score'] = $potentialScore;
        $competencyComponent['problem_solving_score'] = $problemSolvingScore;
        $competencyComponent['score'] = $competencyAggregateScore;

        $subComponents = [
            'education' => $educationComponent['education_points'] ?? 0,
            'experience' => $educationComponent['experience_points'] ?? 0,
            'potential' => $potentialScore,
            'competency_technical' => $competencyComponent['technical_points'] ?? 0,
            'competency_managerial' => $competencyComponent['managerial_points'] ?? 0,
            'competency_communication' => $competencyComponent['communication_points'] ?? 0,
            'competency_total' => $competencyAggregateScore,
            'problem_solving' => $problemSolvingScore,
            'financial_with_authority' => $financialComponent['has_authority']
                ? ($financialComponent['score'] ?? 0)
                : 0,
            'financial_without_authority' => $financialComponent['has_authority']
                ? 0
                : ($financialComponent['score'] ?? 0),
            'asset_management' => $assetComponent['score'] ?? 0,
            'subordinate_control' => $subordinateComponent['score'] ?? 0,
            'work_scope' => $relationshipComponent['score'] ?? 0,
            'work_activity' => $activityComponent['score'] ?? 0,
            'work_risk' => $riskComponent['score'] ?? 0,
        ];

        $baseScore = $educationComponent['score']
            + $competencyAggregateScore
            + $assetComponent['score']
            + $subordinateComponent['score']
            + $relationshipComponent['score']
            + $activityComponent['score']
            + $riskComponent['score'];

        if ($financialComponent['has_authority']) {
            $totals = [
                'with_financial' => $baseScore + ($financialComponent['score'] ?? 0),
                'without_financial' => 0,
            ];
        } else {
            $totals = [
                'with_financial' => 0,
                'without_financial' => $baseScore + ($financialComponent['score'] ?? 0),
            ];
        }

        return [
            'components' => [
                'education_experience' => $educationComponent,
                'potentials' => $potentialComponent,
                'competencies' => $competencyComponent,
                'problem_solving' => $problemSolvingComponent,
                'financial_authority' => $financialComponent,
                'asset_authority' => $assetComponent,
                'subordinate_control' => $subordinateComponent,
                'work_scope' => $relationshipComponent,
                'work_activity' => $activityComponent,
                'work_risk' => $riskComponent,
            ],
            'totals' => $totals,
            'has_financial_authority' => $financialComponent['has_authority'],
            'sub_components' => $subComponents,
        ];
    }

    private function calculateEducationExperience(?JobManagementEducationExperiences $record): array
    {
        $educationLevel = $record?->educationGrade?->level;
        $experienceLevel = $record?->experienceGrade?->level;

        $educationPoints = $this->mapPoints(self::MAP_DEFAULT, $educationLevel);
        $experiencePoints = $this->mapPoints(self::MAP_DEFAULT, $experienceLevel);

        $score = ($educationPoints > 0 && $experiencePoints > 0)
            ? $educationPoints * $experiencePoints
            : 0;

        return [
            'education_level' => $educationLevel,
            'experience_level' => $experienceLevel,
            'education_points' => $educationPoints,
            'experience_points' => $experiencePoints,
            'score' => $score,
        ];
    }

    private function calculatePotentials(Collection $competencies): array
    {
        $psychologicalLevels = $this->collectLevelsByType($competencies, 'Psychological');
        $averageLevel = $psychologicalLevels->count()
            ? (float) ceil($psychologicalLevels->avg())
            : null;

        $points = $averageLevel !== null
            ? $this->mapPoints(self::MAP_DEFAULT, (int) $averageLevel)
            : 0;

        return [
            'average_level' => $averageLevel,
            'items' => $psychologicalLevels->values()->all(),
            'score' => $points,
        ];
    }

    private function calculateCompetencies(Collection $competencies): array
    {
        $technicalLevels = $this->collectLevelsByType($competencies, 'Technical', 'Technical Competency');
        $managerialLevels = $this->collectLevelsByType($competencies, 'Managerial', 'Manajerial');
        $communicationLevel = $this->findSpecificCompetencyLevel(
            $competencies,
            'Communicating & Influencing Skill'
        );

        $technicalAverage = $technicalLevels->count()
            ? (int) ceil($technicalLevels->avg())
            : null;
        $managerialAverage = $managerialLevels->count()
            ? (int) ceil($managerialLevels->avg())
            : null;

        $technicalPoints = $technicalAverage !== null
            ? $this->mapPoints(self::MAP_EXTENDED, $technicalAverage)
            : 0;
        $managerialPoints = $managerialAverage !== null
            ? $this->mapPoints(self::MAP_DEFAULT, $managerialAverage)
            : 0;
        $communicationPoints = $this->mapPoints(
            self::MAP_COMMUNICATION,
            $communicationLevel ?? 1
        );

        $score = ($technicalPoints > 0 && $managerialPoints > 0 && $communicationPoints > 0)
            ? $technicalPoints * $managerialPoints * $communicationPoints
            : 0;

        return [
            'technical_average_level' => $technicalAverage,
            'managerial_average_level' => $managerialAverage,
            'communication_level' => $communicationLevel,
            'technical_points' => $technicalPoints,
            'managerial_points' => $managerialPoints,
            'communication_points' => $communicationPoints,
            'score' => $score,
        ];
    }

    private function calculateProblemSolving(Collection $competencies): array
    {
        $problemItems = $competencies->filter(function ($item) {
            return $item->jobManagementValue?->type === 'Problem Solving & Decision Making';
        });

        $environmentLevel = $problemItems->first(function ($item) {
            return stripos($item->jobManagementValue?->descriptions ?? '', 'Environment') !== false;
        })?->jobManagementValue?->level;

        $challengeLevel = $problemItems->first(function ($item) {
            return stripos($item->jobManagementValue?->descriptions ?? '', 'Chalenge') !== false
                || stripos($item->jobManagementValue?->descriptions ?? '', 'Challenge') !== false;
        })?->jobManagementValue?->level;

        $environmentPoints = $this->mapPoints(self::MAP_EXTENDED, $environmentLevel);
        $challengePoints = $this->mapPoints(self::MAP_DEFAULT, $challengeLevel);

        $score = ($environmentPoints > 0 && $challengePoints > 0)
            ? $environmentPoints * $challengePoints
            : 0;

        return [
            'environment_level' => $environmentLevel,
            'challenge_level' => $challengeLevel,
            'environment_points' => $environmentPoints,
            'challenge_points' => $challengePoints,
            'score' => $score,
        ];
    }

    private function calculateFinancialAuthority(?JobManagementFinancial $record): array
    {
        $hasAuthority = (bool) ($record?->is_authorized);

        $moneyLevel = $record?->cashValue?->level;
        $authorityLevel = $record?->authorityValue?->level;
        $impactLevel = $record?->impactValue?->level;

        $moneyPoints = $hasAuthority
            ? $this->mapPoints(self::MAP_EXTENDED, $moneyLevel)
            : 0;
        $authorityPoints = $this->mapPoints(self::MAP_EXTENDED, $authorityLevel);
        $impactPoints = $this->mapPoints(self::MAP_EXTENDED, $impactLevel);

        if ($hasAuthority) {
            $score = ($moneyPoints > 0 && $authorityPoints > 0 && $impactPoints > 0)
                ? $moneyPoints * $authorityPoints * $impactPoints
                : 0;
        } else {
            $score = ($authorityPoints > 0 && $impactPoints > 0)
                ? $authorityPoints * $impactPoints
                : 0;
        }

        $alternateScore = $score;

        return [
            'has_authority' => $hasAuthority,
            'money_level' => $moneyLevel,
            'authority_level' => $authorityLevel,
            'impact_level' => $impactLevel,
            'money_points' => $moneyPoints,
            'authority_points' => $authorityPoints,
            'impact_points' => $impactPoints,
            'score' => $score,
            'alternate_score' => $alternateScore,
        ];
    }

    private function calculateAssetAuthority(?JobManagementAsset $record): array
    {
        $valuePoints = $this->mapPoints(self::MAP_LINEAR_8, $record?->assetValue?->level);
        $authorityPoints = $this->mapPoints(self::MAP_DEFAULT, $record?->assetAuthority?->level);

        return [
            'asset_value_level' => $record?->assetValue?->level,
            'asset_authority_level' => $record?->assetAuthority?->level,
            'asset_value_points' => $valuePoints,
            'asset_authority_points' => $authorityPoints,
            'score' => ($valuePoints > 0 && $authorityPoints > 0)
                ? $valuePoints * $authorityPoints
                : 0,
        ];
    }

    private function calculateSubordinateControl(?JobManagementSubordinateControl $record): array
    {
        $points = $this->mapPoints(self::MAP_DEFAULT, $record?->value?->level);

        return [
            'level' => $record?->value?->level,
            'points' => $points,
            'score' => $points,
        ];
    }

    private function calculateWorkScope(?JobManagementRelationship $record): array
    {
        $scopePoints = $this->mapPoints(self::MAP_DEFAULT, $record?->scopeValue?->level);
        $frequencyPoints = $this->mapPoints(self::MAP_LINEAR_5, $record?->frequencyValue?->level);

        return [
            'scope_level' => $record?->scopeValue?->level,
            'frequency_level' => $record?->frequencyValue?->level,
            'scope_points' => $scopePoints,
            'frequency_points' => $frequencyPoints,
            'score' => ($scopePoints > 0 && $frequencyPoints > 0)
                ? $scopePoints * $frequencyPoints
                : 0,
        ];
    }

    private function calculateWorkActivity(?JobManagementWorkingActivity $record): array
    {
        $points = $this->mapPoints(self::MAP_DEFAULT, $record?->value?->level);

        return [
            'level' => $record?->value?->level,
            'points' => $points,
            'score' => $points,
        ];
    }

    private function calculateWorkRisk(?JobManagementWorkingRisk $record): array
    {
        $environmentPoints = $this->mapPoints(self::MAP_LINEAR_5, $record?->environmentValue?->level);
        $hazardPoints = $this->mapPoints(self::MAP_LINEAR_5, $record?->hazardValue?->level);

        return [
            'environment_level' => $record?->environmentValue?->level,
            'hazard_level' => $record?->hazardValue?->level,
            'environment_points' => $environmentPoints,
            'hazard_points' => $hazardPoints,
            'score' => ($environmentPoints > 0 && $hazardPoints > 0)
                ? $environmentPoints * $hazardPoints
                : 0,
        ];
    }

    private function collectLevelsByType(Collection $competencies, string $legacyType, ?string $masterField = null): Collection
    {
        return $competencies
            ->map(function ($item) use ($legacyType, $masterField) {
                if ($item->jobManagementValue && $item->jobManagementValue->type === $legacyType) {
                    return $item->jobManagementValue->level;
                }

                if (
                    $masterField
                    && $item->competency
                    && $item->competencyValue
                    && $item->competency->field === $masterField
                ) {
                    return $item->competencyValue->level;
                }

                return null;
            })
            ->filter()
            ->values();
    }

    private function findSpecificCompetencyLevel(Collection $competencies, string $needle): ?int
    {
        $match = $competencies->first(function ($item) use ($needle) {
            if (!$item->jobManagementValue) {
                return false;
            }

            return $item->jobManagementValue->type === $needle
                || stripos($item->jobManagementValue->descriptions ?? '', $needle) !== false;
        });

        return $match?->jobManagementValue?->level;
    }

    private function mapPoints(array $map, $level): int
    {
        if ($level === null) {
            return 0;
        }

        $level = (int) $level;

        if ($level <= 0) {
            return 0;
        }

        if (!array_key_exists($level, $map)) {
            $level = max(array_keys($map));
        }

        return (int) $map[$level];
    }

    private function emptyResult(): array
    {
        return [
            'components' => [],
            'totals' => [
                'with_financial' => 0,
                'without_financial' => 0,
            ],
            'has_financial_authority' => false,
            'sub_components' => [],
        ];
    }

    private function persistResults(array $results): void
    {
        if (empty($results)) {
            return;
        }

        $timestamp = now();

        foreach ($results as $organizationId => $result) {
            $isComplete = $this->isResultComplete($result);
            JobManagementScore::updateOrCreate(
                ['organization_id' => $organizationId],
                [
                    'job_value_with_financial' => $result['totals']['with_financial'] ?? 0,
                    'job_value_without_financial' => $result['totals']['without_financial'] ?? 0,
                    'has_financial_authority' => (bool) ($result['has_financial_authority'] ?? false),
                    'components' => $result['components'] ?? [],
                    'sub_component_points' => $result['sub_components'] ?? [],
                    'calculated_at' => $timestamp,
                    'is_complete' => $isComplete,
                    'completed_at' => $isComplete ? $timestamp : null,
                ]
            );
        }
    }

    private function isResultComplete(array $result): bool
    {
        $sub = $result['sub_components'] ?? [];

        $required = [
            'education',
            'experience',
            'potential',
            'competency_technical',
            'competency_managerial',
            'problem_solving',
            'financial_with_authority',
            'financial_without_authority',
            'asset_management',
            'subordinate_control',
            'work_scope',
            'work_activity',
            'work_risk',
        ];

        foreach ($required as $key) {
            if ($key === 'financial_with_authority' || $key === 'financial_without_authority') {
                continue;
            }

            $value = $sub[$key] ?? 0;
            if ($value <= 0) {
                return false;
            }
        }

        $financialWith = $sub['financial_with_authority'] ?? 0;
        $financialWithout = $sub['financial_without_authority'] ?? 0;

        return $financialWith > 0 || $financialWithout > 0;
    }
}
