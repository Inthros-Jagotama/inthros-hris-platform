<?php
namespace Database\Seeders;

use App\Models\Settings\Company;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;

class Pph21TaxBracketSeeder extends Seeder
{
    public function run(): void
    {
        $companyId = Company::first()->id;

        $pph21ComponentId = DB::table('salary_components')
            ->where('code', 'PPH21_TAX')
            ->where('company_id', $companyId)
            ->value('id');

        // Tax brackets
        $brackets = [
            ['bracket_order' => 1, 'lower_bound' => 0,          'upper_bound' => 60000000,    'rate_percent' => 5],
            ['bracket_order' => 2, 'lower_bound' => 60000000,   'upper_bound' => 250000000,   'rate_percent' => 15],
            ['bracket_order' => 3, 'lower_bound' => 250000000,  'upper_bound' => 500000000,   'rate_percent' => 25],
            ['bracket_order' => 4, 'lower_bound' => 500000000,  'upper_bound' => 5000000000,  'rate_percent' => 30],
            ['bracket_order' => 5, 'lower_bound' => 5000000000, 'upper_bound' => null,         'rate_percent' => 35],
        ];

        foreach ($brackets as $bracket) {
            DB::table('pph21_tax_brackets')->updateOrInsert(
                [
                    'company_id'           => $companyId,
                    'bracket_order'        => $bracket['bracket_order'],
                    'effective_start_date' => '2026-01-01',
                ],
                [
                    'uuid'                 => (string) Str::uuid(),
                    'company_id'           => $companyId,
                    'bracket_order'        => $bracket['bracket_order'],
                    'lower_bound'          => $bracket['lower_bound'],
                    'upper_bound'          => $bracket['upper_bound'],
                    'rate_percent'         => $bracket['rate_percent'],
                    'effective_start_date' => '2026-01-01',
                    'effective_end_date'   => null,
                    'status'               => 'ACTIVE',
                    'created_at'           => now(),
                    'updated_at'           => now(),
                ]
            );
        }

        // PPh21 settings
        DB::table('pph21_settings')->updateOrInsert(
            [
                'company_id'   => $companyId,
                'setting_code' => 'PPH21_2026_REGULAR_GROSS',
            ],
            [
                'uuid'                              => (string) Str::uuid(),
                'company_id'                        => $companyId,
                'setting_code'                      => 'PPH21_2026_REGULAR_GROSS',
                'setting_name'                      => 'PPh21 2026 Regular Gross MVP',
                'calculation_method'                => 'REGULAR_GROSS_ANNUALIZED',
                'default_tax_method'                => 'GROSS',
                'pph21_component_id'                => $pph21ComponentId,
                'occupational_expense_rate_percent' => 5.0,
                'occupational_expense_max_monthly'  => 500000,
                'occupational_expense_max_yearly'   => 6000000,
                'deduct_bpjs_health_employee'       => false,
                'deduct_bpjs_jht_employee'          => true,
                'deduct_bpjs_jp_employee'           => true,
                'annualization_months'              => 12,
                'pkp_rounding_unit'                 => 1000,
                'non_npwp_multiplier_percent'       => 100,
                'rounding_mode'                     => 'ROUND',
                'effective_start_date'              => '2026-01-01',
                'effective_end_date'                => null,
                'status'                            => 'ACTIVE',
                'notes'                             => 'MVP regular gross annualized method. Review tax rules before production use.',
                'created_at'                        => now(),
                'updated_at'                        => now(),
            ]
        );
    }
}
