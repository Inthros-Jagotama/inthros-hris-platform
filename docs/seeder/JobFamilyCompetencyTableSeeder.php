<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use App\Models\JobFamily;
use App\Models\Masters\Competency;

class JobFamilyCompetencyTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $jobFamilies = JobFamily::all()->mapWithKeys(function ($jobFamily) {
            return [trim($jobFamily->name) => $jobFamily->id];
        });

        $technicalCompetencies = Competency::query()
            ->where('field', 'Technical Competency')
            ->whereNotNull('cluster')
            ->get();

        $managerialCompetencies = Competency::query()
            ->where('field', 'Manajerial')
            ->get();

        $payload = [];
        $now = now();

        foreach ($technicalCompetencies as $competency) {
            $clusterName = trim((string) $competency->cluster);
            if (!$clusterName || !$jobFamilies->has($clusterName)) {
                continue;
            }

            $payload[] = [
                'job_family_id' => $jobFamilies->get($clusterName),
                'competency_id' => $competency->id,
                'created_at' => $now,
                'updated_at' => $now,
            ];
        }

        if ($managerialCompetencies->isNotEmpty()) {
            foreach ($jobFamilies as $jobFamilyId) {
                foreach ($managerialCompetencies as $competency) {
                    $payload[] = [
                        'job_family_id' => $jobFamilyId,
                        'competency_id' => $competency->id,
                        'created_at' => $now,
                        'updated_at' => $now,
                    ];
                }
            }
        }

        if (empty($payload)) {
            return;
        }

        foreach (array_chunk($payload, 500) as $chunk) {
            DB::table('job_family_competencies')->upsert(
                $chunk,
                ['job_family_id', 'competency_id'],
                ['updated_at']
            );
        }
    }
}
