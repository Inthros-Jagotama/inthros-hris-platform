<?php

namespace Database\Seeders;

use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;
use App\Models\JobFamily;
use Maatwebsite\Excel\Facades\Excel;
use Illuminate\Support\Facades\DB;

class JobFamilyTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // Path to the Excel file
        $filePath = base_path('Daftar Kompetensi.xlsx');

        if (!file_exists($filePath)) {
            $this->command->error('Excel file not found: ' . $filePath);
            return;
        }

        try {
            // Load the Excel file
            $spreadsheet = \PhpOffice\PhpSpreadsheet\IOFactory::load($filePath);
            $worksheet = $spreadsheet->getActiveSheet();

            // Get the highest row and column
            $highestRow = $worksheet->getHighestRow();
            $highestColumn = $worksheet->getHighestColumn();

            $fieldColumn = 'C'; // Column for "bidang" (field)
            $clusterColumn = 'D'; // Column for "rumpun" (cluster/job family)

            $jobFamilies = [];
            $startRow = 16;

            // Iterate through rows starting from D16
            for ($row = $startRow; $row <= $highestRow; $row++) {
                // Get the field value (column C)
                $fieldValue = $worksheet->getCell($fieldColumn . $row)->getValue();

                // Check if the field is "Technical Competency"
                if ($fieldValue && trim($fieldValue) === 'Technical Competency') {
                    // Get the cluster value (column D) - this will be our job family
                    $clusterValue = $worksheet->getCell($clusterColumn . $row)->getValue();

                    if ($clusterValue && trim($clusterValue) !== '') {
                        $jobFamilyName = trim($clusterValue);

                        // Add to array if not already present
                        if (!in_array($jobFamilyName, $jobFamilies)) {
                            $jobFamilies[] = $jobFamilyName;
                        }
                    }
                }
            }

            // Insert unique job families into the database
            foreach ($jobFamilies as $jobFamilyName) {
                JobFamily::updateOrCreate([
                    'name' => $jobFamilyName
                ]);
            }

            $this->command->info('Successfully seeded ' . count($jobFamilies) . ' job families from Excel file.');
        } catch (\Exception $e) {
            $this->command->error('Error reading Excel file: ' . $e->getMessage());
        }
    }
}
