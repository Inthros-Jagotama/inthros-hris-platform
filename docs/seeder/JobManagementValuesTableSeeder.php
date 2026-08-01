<?php

namespace Database\Seeders;

use App\Models\JobManagement\JobManagementValue;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;

class JobManagementValuesTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $valuesByType = [
            'Pendidikan' => [
                ['level' => 1, 'descriptions' => 'SMA'],
                ['level' => 2, 'descriptions' => 'D3'],
                ['level' => 3, 'descriptions' => 'S1'],
                ['level' => 4, 'descriptions' => 'S2'],
                ['level' => 5, 'descriptions' => 'S3'],
            ],
            'Pengalaman Kerja' => [
                ['level' => 1, 'descriptions' => '0-2 Tahun'],
                ['level' => 2, 'descriptions' => '3-5 Tahun'],
                ['level' => 3, 'descriptions' => '6-8 Tahun'],
                ['level' => 4, 'descriptions' => '9-11 Tahun'],
                ['level' => 5, 'descriptions' => '> 12 Tahun'],
            ],
            'Jurusan' => [
                ['level' => 1, 'descriptions' => 'Teknik Informatika'],
                ['level' => 1, 'descriptions' => 'Sistem Informasi'],
                ['level' => 1, 'descriptions' => 'Teknik Sipil'],
                ['level' => 1, 'descriptions' => 'Manajemen'],
            ],
            'Bidang Pekerjaan' => [
                ['level' => 1, 'descriptions' => 'Sumber Daya Manusia'],
                ['level' => 1, 'descriptions' => 'Operasional'],
                ['level' => 1, 'descriptions' => 'Keuangan'],
                ['level' => 1, 'descriptions' => 'Pemasaran'],
            ],
        ];

        $codeMap = [
            'Psychological' => [
                'Kecerdasan'               => 'kecerdasan',
                'Innovation & Creativity'  => 'innovation_creativity',
                'Self Confidence'          => 'self_confidence',
                'Flexibility'              => 'flexibility',
                'Tenacity'                 => 'tenacity',
                'Continuous Learning'      => 'continuous_learning',
            ],
            'Technical' => [
                'Competency Based Human Resources Management' => 'competency_based_human_resources_management',
                'Competency Development'                      => 'competency_development',
                'People Development'                          => 'people_development',
                'Career Management'                           => 'career_management',
                'HR Assessment'                               => 'hr_assessment',
                'Recruitement & Selection'                    => 'recruitement_selection',
                'Job Analysis & Evaluation'                   => 'job_analysis_evaluation',
                'Organizational Development'                  => 'organizational_development',
                'Human Resources Information System'          => 'human_resources_information_system',
                'Workload Analysis'                           => 'workload_analysis',
                'Performance Apraisal'                        => 'performance_apraisal',
                'Remuneration Manajemen'                      => 'remuneration_manajemen',
                'Reward & Punisment Management'               => 'reward_punisment_management',
                'Health & Safety Environment'                 => 'health_safety_environment',
                'Hubungan Industrial'                         => 'hubungan_industrial',
                'Budgeting'                                   => 'budgeting',
            ],
            'Managerial' => [
                'Integrity'                => 'integrity',
                'Achievement Orientation'  => 'achievement_orientation',
                'Building Partnership'     => 'building_partnership',
                'Planning & Organizing'    => 'planning_organizing',
                'Leadership'               => 'leadership',
                'Developing Others'        => 'developing_others',
            ],
            'Problem Solving & Decision Making' => [
                'Thinking Environment'     => 'thinking_environment',
                'Thinking Chalenge'        => 'thinking_chalenge',
            ],
            'Communicating & Influencing Skill' => [
                'Communicating & Influencing Skill' => 'communicating_influencing_skill',
            ]
        ];

        $additionalValues = [
            'Psychological' => [
                // Kecerdasan - 1
                ['level' => 1, 'descriptions' => 'Kecerdasan', 'note' => 'Kurang'],
                ['level' => 2, 'descriptions' => 'Kecerdasan', 'note' => 'Cukup'],
                ['level' => 3, 'descriptions' => 'Kecerdasan', 'note' => 'Rata-rata'],
                ['level' => 4, 'descriptions' => 'Kecerdasan', 'note' => 'Diatas rata-rata'],
                ['level' => 5, 'descriptions' => 'Kecerdasan', 'note' => 'Istimewa'],

                // Innovation & Creativity - 2
                ['level' => 1, 'descriptions' => 'Innovation & Creativity', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Innovation & Creativity', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Innovation & Creativity', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Innovation & Creativity', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Innovation & Creativity', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Innovation & Creativity', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Innovation & Creativity', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Innovation & Creativity', 'note' => 'Unique Authority'],

                // Self Confidence - 3
                ['level' => 1, 'descriptions' => 'Self Confidence', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Self Confidence', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Self Confidence', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Self Confidence', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Self Confidence', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Self Confidence', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Self Confidence', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Self Confidence', 'note' => 'Unique Authority'],

                // Flexibility - 4
                ['level' => 1, 'descriptions' => 'Flexibility', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Flexibility', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Flexibility', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Flexibility', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Flexibility', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Flexibility', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Flexibility', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Flexibility', 'note' => 'Unique Authority'],

                // Tenacity - 5
                ['level' => 1, 'descriptions' => 'Tenacity', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Tenacity', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Tenacity', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Tenacity', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Tenacity', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Tenacity', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Tenacity', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Tenacity', 'note' => 'Unique Authority'],

                // Continuous Learning - 6

                ['level' => 1, 'descriptions' => 'Continuous Learning', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Continuous Learning', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Continuous Learning', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Continuous Learning', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Continuous Learning', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Continuous Learning', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Continuous Learning', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Continuous Learning', 'note' => 'Unique Authority'],
            ],
            'Technical' => [
                // Competency Based Human Resources Management
                ['level' => 1, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Competency Based Human Resources Management', 'note' => 'Unique Authority'],

                // Competency Development
                ['level' => 1, 'descriptions' => 'Competency Development', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Competency Development', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Competency Development', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Competency Development', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Competency Development', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Competency Development', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Competency Development', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Competency Development', 'note' => 'Unique Authority'],

                // People Development
                ['level' => 1, 'descriptions' => 'People Development', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'People Development', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'People Development', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'People Development', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'People Development', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'People Development', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'People Development', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'People Development', 'note' => 'Unique Authority'],

                // Career Management
                ['level' => 1, 'descriptions' => 'Career Management', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Career Management', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Career Management', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Career Management', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Career Management', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Career Management', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Career Management', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Career Management', 'note' => 'Unique Authority'],

                // HR Assessment
                ['level' => 1, 'descriptions' => 'HR Assessment', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'HR Assessment', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'HR Assessment', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'HR Assessment', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'HR Assessment', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'HR Assessment', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'HR Assessment', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'HR Assessment', 'note' => 'Unique Authority'],

                // Recruitement & Selection
                ['level' => 1, 'descriptions' => 'Recruitement & Selection', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Recruitement & Selection', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Recruitement & Selection', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Recruitement & Selection', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Recruitement & Selection', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Recruitement & Selection', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Recruitement & Selection', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Recruitement & Selection', 'note' => 'Unique Authority'],

                // Job Analysis & Evaluation
                ['level' => 1, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Job Analysis & Evaluation', 'note' => 'Unique Authority'],

                // Organizational Development
                ['level' => 1, 'descriptions' => 'Organizational Development', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Organizational Development', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Organizational Development', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Organizational Development', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Organizational Development', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Organizational Development', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Organizational Development', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Organizational Development', 'note' => 'Unique Authority'],

                // Human Resources Information System
                ['level' => 1, 'descriptions' => 'Human Resources Information System', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Human Resources Information System', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Human Resources Information System', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Human Resources Information System', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Human Resources Information System', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Human Resources Information System', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Human Resources Information System', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Human Resources Information System', 'note' => 'Unique Authority'],

                // Workload Analysis
                ['level' => 1, 'descriptions' => 'Workload Analysis', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Workload Analysis', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Workload Analysis', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Workload Analysis', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Workload Analysis', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Workload Analysis', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Workload Analysis', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Workload Analysis', 'note' => 'Unique Authority'],

                // Performance Apraisal
                ['level' => 1, 'descriptions' => 'Performance Apraisal', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Performance Apraisal', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Performance Apraisal', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Performance Apraisal', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Performance Apraisal', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Performance Apraisal', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Performance Apraisal', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Performance Apraisal', 'note' => 'Unique Authority'],

                // Remuneration Manajemen
                ['level' => 1, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Remuneration Manajemen', 'note' => 'Unique Authority'],

                // Reward & Punisment Management
                ['level' => 1, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Reward & Punisment Management', 'note' => 'Unique Authority'],

                // Health & Safety Environment
                ['level' => 1, 'descriptions' => 'Health & Safety Environment', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Health & Safety Environment', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Health & Safety Environment', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Health & Safety Environment', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Health & Safety Environment', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Health & Safety Environment', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Health & Safety Environment', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Health & Safety Environment', 'note' => 'Unique Authority'],

                // Hubungan Industrial
                ['level' => 1, 'descriptions' => 'Hubungan Industrial', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Hubungan Industrial', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Hubungan Industrial', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Hubungan Industrial', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Hubungan Industrial', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Hubungan Industrial', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Hubungan Industrial', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Hubungan Industrial', 'note' => 'Unique Authority'],

                // Budgeting
                ['level' => 1, 'descriptions' => 'Budgeting', 'note' => 'Primary'],
                ['level' => 2, 'descriptions' => 'Budgeting', 'note' => 'Elementary Vocational'],
                ['level' => 3, 'descriptions' => 'Budgeting', 'note' => 'Vocational'],
                ['level' => 4, 'descriptions' => 'Budgeting', 'note' => 'Advanced Vocational'],
                ['level' => 5, 'descriptions' => 'Budgeting', 'note' => 'Basic Professional'],
                ['level' => 6, 'descriptions' => 'Budgeting', 'note' => 'Seasoned Professional'],
                ['level' => 7, 'descriptions' => 'Budgeting', 'note' => 'Professional Mastery'],
                ['level' => 8, 'descriptions' => 'Budgeting', 'note' => 'Unique Authority'],

            ],
            'Managerial' => [
                // Integrity
                ['level' => 1, 'descriptions' => 'Integrity', 'note' => 'Task'],
                ['level' => 2, 'descriptions' => 'Integrity', 'note' => 'Supervisory'],
                ['level' => 3, 'descriptions' => 'Integrity', 'note' => 'Managerial '],
                ['level' => 4, 'descriptions' => 'Integrity', 'note' => 'Diverse Managerial'],
                ['level' => 5, 'descriptions' => 'Integrity', 'note' => 'Total Managerial'],

                // Achievement Orientation
                ['level' => 1, 'descriptions' => 'Achievement Orientation', 'note' => 'Task'],
                ['level' => 2, 'descriptions' => 'Achievement Orientation', 'note' => 'Supervisory'],
                ['level' => 3, 'descriptions' => 'Achievement Orientation', 'note' => 'Managerial '],
                ['level' => 4, 'descriptions' => 'Achievement Orientation', 'note' => 'Diverse Managerial'],
                ['level' => 5, 'descriptions' => 'Achievement Orientation', 'note' => 'Total Managerial'],

                // Building Partnership
                ['level' => 1, 'descriptions' => 'Building Partnership', 'note' => 'Task'],
                ['level' => 2, 'descriptions' => 'Building Partnership', 'note' => 'Supervisory'],
                ['level' => 3, 'descriptions' => 'Building Partnership', 'note' => 'Managerial '],
                ['level' => 4, 'descriptions' => 'Building Partnership', 'note' => 'Diverse Managerial'],
                ['level' => 5, 'descriptions' => 'Building Partnership', 'note' => 'Total Managerial'],

                // Planning & Organizing
                ['level' => 1, 'descriptions' => 'Planning & Organizing', 'note' => 'Task'],
                ['level' => 2, 'descriptions' => 'Planning & Organizing', 'note' => 'Supervisory'],
                ['level' => 3, 'descriptions' => 'Planning & Organizing', 'note' => 'Managerial '],
                ['level' => 4, 'descriptions' => 'Planning & Organizing', 'note' => 'Diverse Managerial'],
                ['level' => 5, 'descriptions' => 'Planning & Organizing', 'note' => 'Total Managerial'],

                // Leadership
                ['level' => 1, 'descriptions' => 'Leadership', 'note' => 'Task'],
                ['level' => 2, 'descriptions' => 'Leadership', 'note' => 'Supervisory'],
                ['level' => 3, 'descriptions' => 'Leadership', 'note' => 'Managerial '],
                ['level' => 4, 'descriptions' => 'Leadership', 'note' => 'Diverse Managerial'],
                ['level' => 5, 'descriptions' => 'Leadership', 'note' => 'Total Managerial'],

                // Developing Others
                ['level' => 1, 'descriptions' => 'Developing Others', 'note' => 'Task'],
                ['level' => 2, 'descriptions' => 'Developing Others', 'note' => 'Supervisory'],
                ['level' => 3, 'descriptions' => 'Developing Others', 'note' => 'Managerial '],
                ['level' => 4, 'descriptions' => 'Developing Others', 'note' => 'Diverse Managerial'],
                ['level' => 5, 'descriptions' => 'Developing Others', 'note' => 'Total Managerial'],
            ],
            'Problem Solving & Decision Making' => [
                // Thinking Environment
                ['level' => 1, 'descriptions' => 'Thinking Environment', 'note' => 'Berulang-ulang'],
                ['level' => 2, 'descriptions' => 'Thinking Environment', 'note' => 'Bermotif'],
                ['level' => 3, 'descriptions' => 'Thinking Environment', 'note' => 'Variabel'],
                ['level' => 4, 'descriptions' => 'Thinking Environment', 'note' => 'Adaptif'],
                ['level' => 5, 'descriptions' => 'Thinking Environment', 'note' => 'Belum dipetakan'],
                ['level' => 6, 'descriptions' => 'Thinking Environment', 'note' => 'Didefinisikan secara luas'],
                ['level' => 7, 'descriptions' => 'Thinking Environment', 'note' => 'Didefinisikan Secara Umum'],
                ['level' => 8, 'descriptions' => 'Thinking Environment', 'note' => 'Didefinisikan Secara Abstrak'],

                // Thinking Chalenge
                ['level' => 1, 'descriptions' => 'Thinking Chalenge', 'note' => 'Berulang-ulang'],
                ['level' => 2, 'descriptions' => 'Thinking Chalenge', 'note' => 'Bermotif'],
                ['level' => 3, 'descriptions' => 'Thinking Chalenge', 'note' => 'Variabel'],
                ['level' => 4, 'descriptions' => 'Thinking Chalenge', 'note' => 'Adaptif'],
                ['level' => 5, 'descriptions' => 'Thinking Chalenge', 'note' => 'Belum dipetakan'],
            ],
            'Communicating & Influencing Skill' => [
                // Thinking Chalenge
                ['level' => 1, 'descriptions' => 'Communicating & Influencing Skill', 'note' => 'Berkomunikasi'],
                ['level' => 2, 'descriptions' => 'Communicating & Influencing Skill', 'note' => 'Alasan'],
                ['level' => 3, 'descriptions' => 'Communicating & Influencing Skill', 'note' => 'Perubahan Perilaku'],
            ],
            'Jumlah Uang' => [
                ['level' => 1, 'descriptions' => '0 - 500 Jt'],
                ['level' => 2, 'descriptions' => '500 Jt - 1 M'],
                ['level' => 3, 'descriptions' => '1 M - 5 M'],
                ['level' => 4, 'descriptions' => '5 M - 10 M'],
                ['level' => 5, 'descriptions' => '> 10 M'],
            ],
            'Wewenang' => [
                ['level' => 1, 'descriptions' => 'Beroperasi dalam instruksi langsung dan rinci dengan pengawasan yang sangat ketat dan berkelanjutan.', 'note' => 'Memiliki Wewenang'],
                ['level' => 2, 'descriptions' => 'Tunduk pada instruksi dan pekerjaan yang ditetapkan, rutinitas, di bawah pengawasan ketat.', 'note' => 'Memiliki Wewenang'],
                ['level' => 3, 'descriptions' => 'Beroperasi sesuai praktik dan prosedur standar, instruksi kerja umum, dengan pengawasan kemajuan dan hasil.', 'note' => 'Memiliki Wewenang'],
                ['level' => 4, 'descriptions' => 'Beroperasi dalam praktik dan prosedur yang tercakup dalam preseden atau kebijakan yang jelas dan peninjauan hasil akhir.', 'note' => 'Memiliki Wewenang'],
                ['level' => 5, 'descriptions' => 'Tunduk pada praktik dan prosedur luas yang tercakup dalam preseden fungsional dan kebijakan serta arahan manajerial.', 'note' => 'Memiliki Wewenang'],
                ['level' => 6, 'descriptions' => 'Tunduk pada arahan umum dan tujuan kebijakan yang ditetapkan secara luas.', 'note' => 'Memiliki Wewenang'],
                ['level' => 7, 'descriptions' => 'Hanya tunduk pada panduan keseluruhan mengenai tujuan organisasi dan orientasi kebijakan strategis.', 'note' => 'Memiliki Wewenang'],
                ['level' => 8, 'descriptions' => 'Berdasarkan ukuran dan kompleksitas organisasi, hanya tunduk pada panduan yang sangat luas dan orientasi umum terhadap tren bisnis.', 'note' => 'Memiliki Wewenang'],
            ],
            'Dampak pada Hasil Akhir (Memiliki Wewenang Keuangan)' => [
                ['level' => 1, 'descriptions' => 'Penyediaan jasa insidentil untuk penggunaan lainnya.'],
                ['level' => 2, 'descriptions' => 'Penyediaan layanan dukungan informasi/pencatatan atau pengoperasian sederhana peralatan pendukung.'],
                ['level' => 3, 'descriptions' => 'Pengoperasian proses atau peralatan yang berhubungan langsung dengan rantai nilai inti bisnis.'],
                ['level' => 4, 'descriptions' => 'Layanan analitis, diagnostik, konsultasi, atau pengoperasian sistem penting yang sangat kompleks.'],
                ['level' => 5, 'descriptions' => 'Memimpin area aktivitas/tim dalam parameter yang jelas atau memberi nasihat pada tingkat kebijakan.'],
                ['level' => 6, 'descriptions' => 'Menghasilkan operasi multi tim atau memastikan penyampaian program strategis serta kebijakan fungsional.'],
            ],
            'Memiliki Wewenang' => [
                ['level' => 1, 'descriptions' => 'Dikendalikan secara ketat : Beroperasi dalam instruksi langsung dan rinci dengan pengawasan yang sangat ketat dan berkelanjutan.', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 2, 'descriptions' => 'Terkendali : Tunduk pada instruksi dan pekerjaan yang ditetapkan, rutinitas, di bawah pengawasan ketat.', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 3, 'descriptions' => 'Terstandar : Beroperasi sesuai praktik dan prosedur standar, instruksi kerja umum dan pengawasan kemajuan dan hasil.', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 4, 'descriptions' => 'Secara umum diatur : Beroperasi dalam praktik dan prosedur yang tercakup dalam preseden atau kebijakan yang jelas dan peninjauan hasil akhir.', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 5, 'descriptions' => 'Terarah dengan Jelas; Tunduk pada praktik dan prosedur luas yang tercakup dalam preseden fungsional dan kebijakan serta arahan manajerial', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 6, 'descriptions' => 'Dipandu : Hanya tunduk pada panduan keseluruhan mengenai tujuan organisasi secara luas dan orientasi kebijakan strategis.', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 7, 'descriptions' => 'Dipandu Secara Strategis : Berdasarkan ukuran dan kompleksitas organisasi, hanya tunduk pada panduan yang sangat luas dan orientasi umum', 'note' => 'Tidak memiliki Wewenang'],
                ['level' => 8, 'descriptions' => 'Dipandu Secara Strategis : Berdasarkan ukuran dan kompleksitas organisasi, hanya tunduk pada panduan yang sangat luas dan orientasi umum dalam menanggapi tren bisnis.', 'note' => 'Tidak memiliki Wewenang'],
            ],
            'Dampak pada Hasil Akhir (Tidak Memiliki Wewenang Keuangan)' => [
                ['level' => 1, 'descriptions' => 'PAnciliary : Penyediaan jasa insidentil untuk penggunaan lainnya'],
                ['level' => 2, 'descriptions' => 'Suportif : Penyediaan layanan dukungan yang biasanya bersifat informasi dan pencatatan dalam suatu departemen. Atau Pengoperasian atau pemeliharaan sederhana peralatan atau mesin sekunder atau pendukung.'],
                ['level' => 3, 'descriptions' => 'Operasional : Pengoperasian proses dan/atau peralatan yang berhubungan langsung dengan rantai nilai inti bisnis ATAU pengoperasian/pemeliharaan peralatan atau sistem penting yang sangat kompleks dan merupakan inti bisnis'],
                ['level' => 4, 'descriptions' => 'Analitik : Penyediaan layanan khusus yang biasanya bersifat analitis, diagnostik, dan konsultasi.ATAU pengoperasian/pemeliharaan peralatan atau sistem penting dan sangat kompleks yang merupakan inti bisnis'],
                ['level' => 5, 'descriptions' => 'Pemandu : Memimpin area aktivitas yang dapat diidentifikasi seperti tim karyawan atau proyek sehari-hari dalam parameter yang ditentukan dengan baik, ATAU Memberikan nasihat dan bimbingan dalam bidang keahlian pada tingkat pengembangan kerangka kerja/kebijakan.'],
                ['level' => 6, 'descriptions' => 'Berpengaruh : Memberikan hasil dari operasi yang terdiri dari beberapa tim (terkait) yang melakukan beragam aktivitas.ATAU memastikan penyampaian program strategis yang efektif. ATAU Memimpin penyediaan kebijakan dan kerangka fungsional yang memungkinkan kinerja organisasi.'],
            ],
            'Nilai Asset' => [
                ['level' => 1, 'descriptions' => '0 - 1 Jt'],
                ['level' => 2, 'descriptions' => '1 - 10 Jt'],
                ['level' => 3, 'descriptions' => '10 - 50 Jt'],
                ['level' => 4, 'descriptions' => '50 - 100 Jt'],
                ['level' => 5, 'descriptions' => '100 - 250 Jt'],
                ['level' => 6, 'descriptions' => '250 - 500 Jt'],
                ['level' => 7, 'descriptions' => '500 Jt - 1 M'],
                ['level' => 8, 'descriptions' => '> 1 M'],
            ],
            'Wewenang Asset' => [
                ['level' => 1, 'descriptions' => 'Tidak ada Aset'],
                ['level' => 2, 'descriptions' => 'Menggunakan Aset'],
                ['level' => 3, 'descriptions' => 'Mengelola Aset'],
                ['level' => 4, 'descriptions' => 'Memeriksa Aset'],
                ['level' => 5, 'descriptions' => 'Memverifikasi Aset'],
                ['level' => 6, 'descriptions' => 'Menyetujui Aset'],
            ],
            'Total Bawahan' => [
                ['level' => 1, 'descriptions' => 'Very Small'],
                ['level' => 2, 'descriptions' => 'Small'],
                ['level' => 3, 'descriptions' => 'Medium'],
                ['level' => 4, 'descriptions' => 'Large'],
                ['level' => 5, 'descriptions' => 'Very Large'],
            ],
            'Aktifitas Fisik' => [
                ['level' => 1, 'descriptions' => 'Banyak duduk sedikit bergerak'],
                ['level' => 2, 'descriptions' => 'Seimbang duduk dan berdiri'],
                ['level' => 3, 'descriptions' => 'Sedikit duduk, banyak berdiri dan berjalan'],
                ['level' => 4, 'descriptions' => 'Aktivitas fisik tinggi, gunakan organ dan indra'],
                ['level' => 5, 'descriptions' => 'Aktivitas sangat tinggi, melakukan pengawasan'],
            ],
            'Lingkungan Kerja' => [
                ['level' => 1, 'descriptions' => 'Tenang, nyaman'],
                ['level' => 2, 'descriptions' => 'Cukup bising dan sibuk'],
                ['level' => 3, 'descriptions' => 'Bising, sibuk dan banyak gangguan'],
                ['level' => 4, 'descriptions' => 'Banyak tantangan dan tekanan kerja'],
                ['level' => 5, 'descriptions' => 'Sangat menekan dan menegangkan'],
            ],
            'Resiko/Bahaya' => [
                ['level' => 1, 'descriptions' => 'Risiko minimum dan bebas bahaya'],
                ['level' => 2, 'descriptions' => 'Risiko kecil dengan sedikit ancaman bahaya'],
                ['level' => 3, 'descriptions' => 'Risiko dan bahaya yang cukup besar'],
                ['level' => 4, 'descriptions' => 'Risiko dan bahaya tinggi'],
                ['level' => 5, 'descriptions' => 'Ancaman bahaya besar yang mematikan'],
            ],
            'Lingkup Hubungan Kerja' => [
                ['level' => 1, 'descriptions' => 'Unit Kerja'],
                ['level' => 2, 'descriptions' => 'Antar Unit Kerja'],
                ['level' => 3, 'descriptions' => 'Lingkup Internal secara Nasional'],
                ['level' => 4, 'descriptions' => 'Lingkup Eksternal secara Nasional'],
                ['level' => 5, 'descriptions' => 'Lingkup Internasional'],
            ],
            'Frekuensi Hubungan Kerja' => [
                ['level' => 1, 'descriptions' => 'Sesekali'],
                ['level' => 2, 'descriptions' => 'Kadang - kadang'],
                ['level' => 3, 'descriptions' => 'Cukup Sering'],
                ['level' => 4, 'descriptions' => 'Sering'],
                ['level' => 5, 'descriptions' => 'Sangat Sering'],
            ],
        ];

        $valuesByType = array_merge($valuesByType, $additionalValues);

        $records = collect($valuesByType)->flatMap(function (array $items, string $type) use ($codeMap) {
            // Get unique descriptions in order of appearance for sorting
            $uniqueDescriptions = [];
            $sortCounter = 1;

            return collect($items)->map(function (array $item) use ($type, $codeMap, &$uniqueDescriptions, &$sortCounter) {
                $desc = $item['descriptions'];

                $code = $codeMap[$type][$desc] ?? Str::of($desc)
                    ->lower()
                    ->replace(['&', '/'], ' ')
                    ->replaceMatches('/[^a-z0-9]+/i', '_')
                    ->trim('_')
                    ->toString();

                // Assign sort value for specific types based on description order
                $sort = null;
                if (in_array($type, ['Psychological', 'Technical', 'Managerial', 'Problem Solving & Decision Making', 'Communicating & Influencing Skill'])) {
                    if (!isset($uniqueDescriptions[$desc])) {
                        $uniqueDescriptions[$desc] = $sortCounter++;
                    }
                    $sort = $uniqueDescriptions[$desc];
                }

                return [
                    'type'                           => $type,
                    'code'                           => $code ?? null,
                    'level'                          => (int) $item['level'],
                    'descriptions'                   => $desc,
                    'note'                           => $item['note'] ?? null,
                    'sort'                           => $sort,
                    'job_management_title_sub_id'    => null,
                    'job_management_title_sub_name'  => null,
                    'created_at'                     => now(),
                    'updated_at'                     => now(),
                ];
            });
        });

        $uniqueRecords = $records
            ->unique(fn($r) => $r['type'] . '|' . $r['code'] . '|' . $r['level'])
            ->values();

        $types = $uniqueRecords->pluck('type')->unique()->values()->all();

        if (!empty($types)) {
            DB::table('job_management_values')->whereIn('type', $types)->delete();
        }

        JobManagementValue::insert($uniqueRecords->all());
    }
}
