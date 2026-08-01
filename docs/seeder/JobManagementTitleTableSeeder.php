<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use App\Models\JobManagement\JobManagementTitle;

class JobManagementTitleTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $jobManagementTitles = [
            [
                'id' => 1,
                'name' => 'IDENTIFIKASI JABATAN',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 2,
                'name' => 'TUJUAN JABATAN',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 3,
                'name' => 'TANGGUNG JAWAB',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 4,
                'name' => 'PENDIDIKAN DAN PENGALAMAN',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 5,
                'name' => 'POTENSI DAN KOMPETENSI YANG HARUS DIMILIKI',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 6,
                'name' => 'MEMILIKI WEWENANG PENGELOLAAN KEUANGAN',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 7,
                'name' => 'TIDAK MEMILIKI WEWENANG PENGELOLAAN KEUANGAN',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 8,
                'name' => 'WEWENANG PENGELOLAAN ASSET',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 9,
                'name' => 'RENTANG KENDALI BAWAHAN',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 10,
                'name' => 'HUBUNGAN KERJA',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 11,
                'name' => 'AKTIVITAS KERJA',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 12,
                'name' => 'RISIKO KERJA',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 13,
                'name' => 'WEWENANG BIDANG SDM',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
            [
                'id' => 14,
                'name' => 'WEWENANG BIDANG OPERASI',
                'descriptions' => null,
                'created_by' => null,
                'updated_by' => null,
                'status' => 1,
            ],
        ];

        foreach ($jobManagementTitles as $jobManagementTitle) {
            JobManagementTitle::updateOrCreate(['id' => $jobManagementTitle['id']], $jobManagementTitle);
        }
    }
}
