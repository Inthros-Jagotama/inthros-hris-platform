<?php

namespace Database\Seeders;

use App\Models\Settings\Ter;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class TerSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $group_a = [
            ['bruto_min' => 0, 'bruto_max' => 5400000, 'rate' => 0.00],
            ['bruto_min' => 5400000, 'bruto_max' => 5650000, 'rate' => 0.25],
            ['bruto_min' => 5650000, 'bruto_max' => 5950000, 'rate' => 0.50],
            ['bruto_min' => 5950000, 'bruto_max' => 6300000, 'rate' => 0.75],
            ['bruto_min' => 6300000, 'bruto_max' => 6750000, 'rate' => 1.00],
            ['bruto_min' => 6750000, 'bruto_max' => 7500000, 'rate' => 1.25],
            ['bruto_min' => 7500000, 'bruto_max' => 8550000, 'rate' => 1.50],
            ['bruto_min' => 8550000, 'bruto_max' => 9650000, 'rate' => 1.75],
            ['bruto_min' => 9650000, 'bruto_max' => 10050000, 'rate' => 2.00],
            ['bruto_min' => 10050000, 'bruto_max' => 10350000, 'rate' => 2.25],
            ['bruto_min' => 10350000, 'bruto_max' => 10700000, 'rate' => 2.50],
            ['bruto_min' => 10700000, 'bruto_max' => 11050000, 'rate' => 3.00],
            ['bruto_min' => 11050000, 'bruto_max' => 11600000, 'rate' => 3.50],
            ['bruto_min' => 11600000, 'bruto_max' => 12500000, 'rate' => 4.00],
            ['bruto_min' => 12500000, 'bruto_max' => 13750000, 'rate' => 5.00],
            ['bruto_min' => 13750000, 'bruto_max' => 15100000, 'rate' => 6.00],
            ['bruto_min' => 15100000, 'bruto_max' => 16950000, 'rate' => 7.00],
            ['bruto_min' => 16950000, 'bruto_max' => 19750000, 'rate' => 8.00],
            ['bruto_min' => 19750000, 'bruto_max' => 24150000, 'rate' => 9.00],
            ['bruto_min' => 24150000, 'bruto_max' => 26450000, 'rate' => 10.00],
            ['bruto_min' => 26450000, 'bruto_max' => 28000000, 'rate' => 11.00],
            ['bruto_min' => 28000000, 'bruto_max' => 30050000, 'rate' => 12.00],
            ['bruto_min' => 30050000, 'bruto_max' => 32400000, 'rate' => 13.00],
            ['bruto_min' => 32400000, 'bruto_max' => 35400000, 'rate' => 14.00],
            ['bruto_min' => 35400000, 'bruto_max' => 39100000, 'rate' => 15.00],
            ['bruto_min' => 39100000, 'bruto_max' => 43850000, 'rate' => 16.00],
            ['bruto_min' => 43850000, 'bruto_max' => 47800000, 'rate' => 17.00],
            ['bruto_min' => 47800000, 'bruto_max' => 51400000, 'rate' => 18.00],
            ['bruto_min' => 51400000, 'bruto_max' => 56300000, 'rate' => 19.00],
            ['bruto_min' => 56300000, 'bruto_max' => 62200000, 'rate' => 20.00],
            ['bruto_min' => 62200000, 'bruto_max' => 68600000, 'rate' => 21.00],
            ['bruto_min' => 68600000, 'bruto_max' => 77500000, 'rate' => 22.00],
            ['bruto_min' => 77500000, 'bruto_max' => 89000000, 'rate' => 23.00],
            ['bruto_min' => 89000000, 'bruto_max' => 103000000, 'rate' => 24.00],
            ['bruto_min' => 103000000, 'bruto_max' => 125000000, 'rate' => 25.00],
            ['bruto_min' => 125000000, 'bruto_max' => 157000000, 'rate' => 26.00],
            ['bruto_min' => 157000000, 'bruto_max' => 206000000, 'rate' => 27.00],
            ['bruto_min' => 206000000, 'bruto_max' => 337000000, 'rate' => 28.00],
            ['bruto_min' => 337000000, 'bruto_max' => 454000000, 'rate' => 29.00],
            ['bruto_min' => 454000000, 'bruto_max' => 550000000, 'rate' => 30.00],
            ['bruto_min' => 550000000, 'bruto_max' => 695000000, 'rate' => 31.00],
            ['bruto_min' => 695000000, 'bruto_max' => 910000000, 'rate' => 32.00],
            ['bruto_min' => 910000000, 'bruto_max' => 1400000000, 'rate' => 33.00],
            ['bruto_min' => 1400000000, 'bruto_max' => null, 'rate' => 34.00],
        ];

        foreach ($group_a as $item) {
            Ter::updateOrInsert([
                'group' => 'A',
                'bruto_min' => $item['bruto_min'],
                'bruto_max' => $item['bruto_max'],
                'rate' => $item['rate'],
            ]);
        }

        $group_b = [
            ['bruto_min' => 0, 'bruto_max' => 6200000, 'rate' => 0.00],
            ['bruto_min' => 6200000, 'bruto_max' => 6500000, 'rate' => 0.25],
            ['bruto_min' => 6500000, 'bruto_max' => 6850000, 'rate' => 0.50],
            ['bruto_min' => 6850000, 'bruto_max' => 7300000, 'rate' => 0.75],
            ['bruto_min' => 7300000, 'bruto_max' => 9200000, 'rate' => 1.00],
            ['bruto_min' => 9200000, 'bruto_max' => 10750000, 'rate' => 1.50],
            ['bruto_min' => 10750000, 'bruto_max' => 11250000, 'rate' => 2.00],
            ['bruto_min' => 11250000, 'bruto_max' => 11600000, 'rate' => 2.50],
            ['bruto_min' => 11600000, 'bruto_max' => 12600000, 'rate' => 3.00],
            ['bruto_min' => 12600000, 'bruto_max' => 13600000, 'rate' => 4.00],
            ['bruto_min' => 13600000, 'bruto_max' => 14950000, 'rate' => 5.00],
            ['bruto_min' => 14950000, 'bruto_max' => 16400000, 'rate' => 6.00],
            ['bruto_min' => 16400000, 'bruto_max' => 18450000, 'rate' => 7.00],
            ['bruto_min' => 18450000, 'bruto_max' => 21850000, 'rate' => 8.00],
            ['bruto_min' => 21850000, 'bruto_max' => 26000000, 'rate' => 9.00],
            ['bruto_min' => 26000000, 'bruto_max' => 27700000, 'rate' => 10.00],
            ['bruto_min' => 27700000, 'bruto_max' => 29350000, 'rate' => 11.00],
            ['bruto_min' => 29350000, 'bruto_max' => 31450000, 'rate' => 12.00],
            ['bruto_min' => 31450000, 'bruto_max' => 33950000, 'rate' => 13.00],
            ['bruto_min' => 33950000, 'bruto_max' => 37100000, 'rate' => 14.00],
            ['bruto_min' => 37100000, 'bruto_max' => 41100000, 'rate' => 15.00],
            ['bruto_min' => 41100000, 'bruto_max' => 45800000, 'rate' => 16.00],
            ['bruto_min' => 45800000, 'bruto_max' => 49500000, 'rate' => 17.00],
            ['bruto_min' => 49500000, 'bruto_max' => 53800000, 'rate' => 18.00],
            ['bruto_min' => 53800000, 'bruto_max' => 58500000, 'rate' => 19.00],
            ['bruto_min' => 58500000, 'bruto_max' => 64000000, 'rate' => 20.00],
            ['bruto_min' => 64000000, 'bruto_max' => 71000000, 'rate' => 21.00],
            ['bruto_min' => 71000000, 'bruto_max' => 80000000, 'rate' => 22.00],
            ['bruto_min' => 80000000, 'bruto_max' => 93000000, 'rate' => 23.00],
            ['bruto_min' => 93000000, 'bruto_max' => 109000000, 'rate' => 24.00],
            ['bruto_min' => 109000000, 'bruto_max' => 129000000, 'rate' => 25.00],
            ['bruto_min' => 129000000, 'bruto_max' => 163000000, 'rate' => 26.00],
            ['bruto_min' => 163000000, 'bruto_max' => 211000000, 'rate' => 27.00],
            ['bruto_min' => 211000000, 'bruto_max' => 374000000, 'rate' => 28.00],
            ['bruto_min' => 374000000, 'bruto_max' => 459000000, 'rate' => 29.00],
            ['bruto_min' => 459000000, 'bruto_max' => 555000000, 'rate' => 30.00],
            ['bruto_min' => 555000000, 'bruto_max' => 704000000, 'rate' => 31.00],
            ['bruto_min' => 704000000, 'bruto_max' => 957000000, 'rate' => 32.00],
            ['bruto_min' => 957000000, 'bruto_max' => 1405000000, 'rate' => 33.00],
            ['bruto_min' => 1405000000, 'bruto_max' => null, 'rate' => 34.00],
        ];

        foreach ($group_b as $item) {
            Ter::updateOrInsert([
                'group' => 'B',
                'bruto_min' => $item['bruto_min'],
                'bruto_max' => $item['bruto_max'],
                'rate' => $item['rate'],
            ]);
        }

        $group_c = [
            ['bruto_min' => 6600000, 'bruto_max' => 6950000, 'rate' => 0.00],
            ['bruto_min' => 6950000, 'bruto_max' => 7350000, 'rate' => 0.25],
            ['bruto_min' => 7350000, 'bruto_max' => 7800000, 'rate' => 0.50],
            ['bruto_min' => 7800000, 'bruto_max' => 8850000, 'rate' => 0.75],
            ['bruto_min' => 8850000, 'bruto_max' => 9800000, 'rate' => 1.00],
            ['bruto_min' => 9800000, 'bruto_max' => 10950000, 'rate' => 1.25],
            ['bruto_min' => 10950000, 'bruto_max' => 11200000, 'rate' => 1.50],
            ['bruto_min' => 11200000, 'bruto_max' => 12050000, 'rate' => 1.75],
            ['bruto_min' => 12050000, 'bruto_max' => 12950000, 'rate' => 2.00],
            ['bruto_min' => 12950000, 'bruto_max' => 14150000, 'rate' => 3.00],
            ['bruto_min' => 14150000, 'bruto_max' => 15550000, 'rate' => 4.00],
            ['bruto_min' => 15550000, 'bruto_max' => 17050000, 'rate' => 5.00],
            ['bruto_min' => 17050000, 'bruto_max' => 19500000, 'rate' => 6.00],
            ['bruto_min' => 19500000, 'bruto_max' => 22700000, 'rate' => 7.00],
            ['bruto_min' => 22700000, 'bruto_max' => 26600000, 'rate' => 8.00],
            ['bruto_min' => 26600000, 'bruto_max' => 28100000, 'rate' => 9.00],
            ['bruto_min' => 28100000, 'bruto_max' => 30100000, 'rate' => 10.00],
            ['bruto_min' => 30100000, 'bruto_max' => 32600000, 'rate' => 11.00],
            ['bruto_min' => 32600000, 'bruto_max' => 35400000, 'rate' => 12.00],
            ['bruto_min' => 35400000, 'bruto_max' => 38900000, 'rate' => 13.00],
            ['bruto_min' => 38900000, 'bruto_max' => 43000000, 'rate' => 14.00],
            ['bruto_min' => 43000000, 'bruto_max' => 47400000, 'rate' => 15.00],
            ['bruto_min' => 47400000, 'bruto_max' => 51200000, 'rate' => 16.00],
            ['bruto_min' => 51200000, 'bruto_max' => 55800000, 'rate' => 17.00],
            ['bruto_min' => 55800000, 'bruto_max' => 60400000, 'rate' => 18.00],
            ['bruto_min' => 60400000, 'bruto_max' => 66700000, 'rate' => 19.00],
            ['bruto_min' => 66700000, 'bruto_max' => 74500000, 'rate' => 20.00],
            ['bruto_min' => 74500000, 'bruto_max' => 83200000, 'rate' => 21.00],
            ['bruto_min' => 83200000, 'bruto_max' => 95600000, 'rate' => 22.00],
            ['bruto_min' => 95600000, 'bruto_max' => 110000000, 'rate' => 23.00],
            ['bruto_min' => 110000000, 'bruto_max' => 134000000, 'rate' => 24.00],
            ['bruto_min' => 134000000, 'bruto_max' => 169000000, 'rate' => 25.00],
            ['bruto_min' => 169000000, 'bruto_max' => 221000000, 'rate' => 26.00],
            ['bruto_min' => 221000000, 'bruto_max' => 390000000, 'rate' => 27.00],
            ['bruto_min' => 390000000, 'bruto_max' => 463000000, 'rate' => 28.00],
            ['bruto_min' => 463000000, 'bruto_max' => 561000000, 'rate' => 29.00],
            ['bruto_min' => 561000000, 'bruto_max' => 709000000, 'rate' => 30.00],
            ['bruto_min' => 709000000, 'bruto_max' => 965000000, 'rate' => 31.00],
            ['bruto_min' => 965000000, 'bruto_max' => 1419000000, 'rate' => 32.00],
            ['bruto_min' => 1419000000, 'bruto_max' => null, 'rate' => 33.00],
        ];

        foreach ($group_c as $item) {
            Ter::updateOrInsert([
                'group' => 'C',
                'bruto_min' => $item['bruto_min'],
                'bruto_max' => $item['bruto_max'],
                'rate' => $item['rate'],
            ]);
        }

    }
}
