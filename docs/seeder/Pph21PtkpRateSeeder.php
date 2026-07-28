<?php
namespace Database\Seeders;

use App\Models\Settings\Company;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;

class Pph21PtkpRateSeeder extends Seeder
{
    public function run(): void
    {
        $companyId = Company::first()->id;

        $ptkpRates = [
            ['ptkp_status' => 'TK/0', 'description' => 'Tidak kawin, 0 tanggungan',                           'annual_amount' => 54000000],
            ['ptkp_status' => 'TK/1', 'description' => 'Tidak kawin, 1 tanggungan',                           'annual_amount' => 58500000],
            ['ptkp_status' => 'TK/2', 'description' => 'Tidak kawin, 2 tanggungan',                           'annual_amount' => 63000000],
            ['ptkp_status' => 'TK/3', 'description' => 'Tidak kawin, 3 tanggungan',                           'annual_amount' => 67500000],
            ['ptkp_status' => 'K/0',  'description' => 'Kawin, 0 tanggungan',                                 'annual_amount' => 58500000],
            ['ptkp_status' => 'K/1',  'description' => 'Kawin, 1 tanggungan',                                 'annual_amount' => 63000000],
            ['ptkp_status' => 'K/2',  'description' => 'Kawin, 2 tanggungan',                                 'annual_amount' => 67500000],
            ['ptkp_status' => 'K/3',  'description' => 'Kawin, 3 tanggungan',                                 'annual_amount' => 72000000],
            ['ptkp_status' => 'K/I/0','description' => 'Kawin, penghasilan istri digabung, 0 tanggungan',     'annual_amount' => 112500000],
            ['ptkp_status' => 'K/I/1','description' => 'Kawin, penghasilan istri digabung, 1 tanggungan',     'annual_amount' => 117000000],
            ['ptkp_status' => 'K/I/2','description' => 'Kawin, penghasilan istri digabung, 2 tanggungan',     'annual_amount' => 121500000],
            ['ptkp_status' => 'K/I/3','description' => 'Kawin, penghasilan istri digabung, 3 tanggungan',     'annual_amount' => 126000000],
        ];

        foreach ($ptkpRates as $rate) {
            DB::table('pph21_ptkp_rates')->updateOrInsert(
                [
                    'company_id'           => $companyId,
                    'ptkp_status'          => $rate['ptkp_status'],
                    'effective_start_date' => '2026-01-01',
                ],
                [
                    'uuid'                 => (string) Str::uuid(),
                    'company_id'           => $companyId,
                    'ptkp_status'          => $rate['ptkp_status'],
                    'description'          => $rate['description'],
                    'annual_amount'        => $rate['annual_amount'],
                    'effective_start_date' => '2026-01-01',
                    'effective_end_date'   => null,
                    'status'               => 'ACTIVE',
                    'created_at'           => now(),
                    'updated_at'           => now(),
                ]
            );
        }
    }
}
