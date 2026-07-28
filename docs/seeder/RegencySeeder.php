<?php

namespace Database\Seeders;

use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class RegencySeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        \DB::table('regencies')->delete();

        \DB::table('regencies')->insert(array(
            0 =>
            array(
                'id' => '1101',
                'province_id' => '11',
                'name' => 'KAB. ACEH SELATAN',
            ),
            1 =>
            array(
                'id' => '1102',
                'province_id' => '11',
                'name' => 'KAB. ACEH TENGGARA',
            ),
            2 =>
            array(
                'id' => '1103',
                'province_id' => '11',
                'name' => 'KAB. ACEH TIMUR',
            ),
            3 =>
            array(
                'id' => '1104',
                'province_id' => '11',
                'name' => 'KAB. ACEH TENGAH',
            ),
            4 =>
            array(
                'id' => '1105',
                'province_id' => '11',
                'name' => 'KAB. ACEH BARAT',
            ),
            5 =>
            array(
                'id' => '1106',
                'province_id' => '11',
                'name' => 'KAB. ACEH BESAR',
            ),
            6 =>
            array(
                'id' => '1107',
                'province_id' => '11',
                'name' => 'KAB. PIDIE',
            ),
            7 =>
            array(
                'id' => '1108',
                'province_id' => '11',
                'name' => 'KAB. ACEH UTARA',
            ),
            8 =>
            array(
                'id' => '1109',
                'province_id' => '11',
                'name' => 'KAB. SIMEULUE',
            ),
            9 =>
            array(
                'id' => '1110',
                'province_id' => '11',
                'name' => 'KAB. ACEH SINGKIL',
            ),
            10 =>
            array(
                'id' => '1111',
                'province_id' => '11',
                'name' => 'KAB. BIREUEN',
            ),
            11 =>
            array(
                'id' => '1112',
                'province_id' => '11',
                'name' => 'KAB. ACEH BARAT DAYA',
            ),
            12 =>
            array(
                'id' => '1113',
                'province_id' => '11',
                'name' => 'KAB. GAYO LUES',
            ),
            13 =>
            array(
                'id' => '1114',
                'province_id' => '11',
                'name' => 'KAB. ACEH JAYA',
            ),
            14 =>
            array(
                'id' => '1115',
                'province_id' => '11',
                'name' => 'KAB. NAGAN RAYA',
            ),
            15 =>
            array(
                'id' => '1116',
                'province_id' => '11',
                'name' => 'KAB. ACEH TAMIANG',
            ),
            16 =>
            array(
                'id' => '1117',
                'province_id' => '11',
                'name' => 'KAB. BENER MERIAH',
            ),
            17 =>
            array(
                'id' => '1118',
                'province_id' => '11',
                'name' => 'KAB. PIDIE JAYA',
            ),
            18 =>
            array(
                'id' => '1171',
                'province_id' => '11',
                'name' => 'KOTA BANDA ACEH',
            ),
            19 =>
            array(
                'id' => '1172',
                'province_id' => '11',
                'name' => 'KOTA SABANG',
            ),
            20 =>
            array(
                'id' => '1173',
                'province_id' => '11',
                'name' => 'KOTA LHOKSEUMAWE',
            ),
            21 =>
            array(
                'id' => '1174',
                'province_id' => '11',
                'name' => 'KOTA LANGSA',
            ),
            22 =>
            array(
                'id' => '1175',
                'province_id' => '11',
                'name' => 'KOTA SUBULUSSALAM',
            ),
            23 =>
            array(
                'id' => '1201',
                'province_id' => '12',
                'name' => 'KAB. TAPANULI TENGAH',
            ),
            24 =>
            array(
                'id' => '1202',
                'province_id' => '12',
                'name' => 'KAB. TAPANULI UTARA',
            ),
            25 =>
            array(
                'id' => '1203',
                'province_id' => '12',
                'name' => 'KAB. TAPANULI SELATAN',
            ),
            26 =>
            array(
                'id' => '1204',
                'province_id' => '12',
                'name' => 'KAB. NIAS',
            ),
            27 =>
            array(
                'id' => '1205',
                'province_id' => '12',
                'name' => 'KAB. LANGKAT',
            ),
            28 =>
            array(
                'id' => '1206',
                'province_id' => '12',
                'name' => 'KAB. KARO',
            ),
            29 =>
            array(
                'id' => '1207',
                'province_id' => '12',
                'name' => 'KAB. DELI SERDANG',
            ),
            30 =>
            array(
                'id' => '1208',
                'province_id' => '12',
                'name' => 'KAB. SIMALUNGUN',
            ),
            31 =>
            array(
                'id' => '1209',
                'province_id' => '12',
                'name' => 'KAB. ASAHAN',
            ),
            32 =>
            array(
                'id' => '1210',
                'province_id' => '12',
                'name' => 'KAB. LABUHANBATU',
            ),
            33 =>
            array(
                'id' => '1211',
                'province_id' => '12',
                'name' => 'KAB. DAIRI',
            ),
            34 =>
            array(
                'id' => '1212',
                'province_id' => '12',
                'name' => 'KAB. TOBA SAMOSIR',
            ),
            35 =>
            array(
                'id' => '1213',
                'province_id' => '12',
                'name' => 'KAB. MANDAILING NATAL',
            ),
            36 =>
            array(
                'id' => '1214',
                'province_id' => '12',
                'name' => 'KAB. NIAS SELATAN',
            ),
            37 =>
            array(
                'id' => '1215',
                'province_id' => '12',
                'name' => 'KAB. PAKPAK BHARAT',
            ),
            38 =>
            array(
                'id' => '1216',
                'province_id' => '12',
                'name' => 'KAB. HUMBANG HASUNDUTAN',
            ),
            39 =>
            array(
                'id' => '1217',
                'province_id' => '12',
                'name' => 'KAB. SAMOSIR',
            ),
            40 =>
            array(
                'id' => '1218',
                'province_id' => '12',
                'name' => 'KAB. SERDANG BEDAGAI',
            ),
            41 =>
            array(
                'id' => '1219',
                'province_id' => '12',
                'name' => 'KAB. BATU BARA',
            ),
            42 =>
            array(
                'id' => '1220',
                'province_id' => '12',
                'name' => 'KAB. PADANG LAWAS UTARA',
            ),
            43 =>
            array(
                'id' => '1221',
                'province_id' => '12',
                'name' => 'KAB. PADANG LAWAS',
            ),
            44 =>
            array(
                'id' => '1222',
                'province_id' => '12',
                'name' => 'KAB. LABUHANBATU SELATAN',
            ),
            45 =>
            array(
                'id' => '1223',
                'province_id' => '12',
                'name' => 'KAB. LABUHANBATU UTARA',
            ),
            46 =>
            array(
                'id' => '1224',
                'province_id' => '12',
                'name' => 'KAB. NIAS UTARA',
            ),
            47 =>
            array(
                'id' => '1225',
                'province_id' => '12',
                'name' => 'KAB. NIAS BARAT',
            ),
            48 =>
            array(
                'id' => '1271',
                'province_id' => '12',
                'name' => 'KOTA MEDAN',
            ),
            49 =>
            array(
                'id' => '1272',
                'province_id' => '12',
                'name' => 'KOTA PEMATANGSIANTAR',
            ),
            50 =>
            array(
                'id' => '1273',
                'province_id' => '12',
                'name' => 'KOTA SIBOLGA',
            ),
            51 =>
            array(
                'id' => '1274',
                'province_id' => '12',
                'name' => 'KOTA TANJUNG BALAI',
            ),
            52 =>
            array(
                'id' => '1275',
                'province_id' => '12',
                'name' => 'KOTA BINJAI',
            ),
            53 =>
            array(
                'id' => '1276',
                'province_id' => '12',
                'name' => 'KOTA TEBING TINGGI',
            ),
            54 =>
            array(
                'id' => '1277',
                'province_id' => '12',
                'name' => 'KOTA PADANG SIDEMPUAN',
            ),
            55 =>
            array(
                'id' => '1278',
                'province_id' => '12',
                'name' => 'KOTA GUNUNGSITOLI',
            ),
            56 =>
            array(
                'id' => '1301',
                'province_id' => '13',
                'name' => 'KAB. PESISIR SELATAN',
            ),
            57 =>
            array(
                'id' => '1302',
                'province_id' => '13',
                'name' => 'KAB. SOLOK',
            ),
            58 =>
            array(
                'id' => '1303',
                'province_id' => '13',
                'name' => 'KAB. SIJUNJUNG',
            ),
            59 =>
            array(
                'id' => '1304',
                'province_id' => '13',
                'name' => 'KAB. TANAH DATAR',
            ),
            60 =>
            array(
                'id' => '1305',
                'province_id' => '13',
                'name' => 'KAB. PADANG PARIAMAN',
            ),
            61 =>
            array(
                'id' => '1306',
                'province_id' => '13',
                'name' => 'KAB. AGAM',
            ),
            62 =>
            array(
                'id' => '1307',
                'province_id' => '13',
                'name' => 'KAB. LIMA PULUH KOTA',
            ),
            63 =>
            array(
                'id' => '1308',
                'province_id' => '13',
                'name' => 'KAB. PASAMAN',
            ),
            64 =>
            array(
                'id' => '1309',
                'province_id' => '13',
                'name' => 'KAB. KEPULAUAN MENTAWAI',
            ),
            65 =>
            array(
                'id' => '1310',
                'province_id' => '13',
                'name' => 'KAB. DHARMASRAYA',
            ),
            66 =>
            array(
                'id' => '1311',
                'province_id' => '13',
                'name' => 'KAB. SOLOK SELATAN',
            ),
            67 =>
            array(
                'id' => '1312',
                'province_id' => '13',
                'name' => 'KAB. PASAMAN BARAT',
            ),
            68 =>
            array(
                'id' => '1371',
                'province_id' => '13',
                'name' => 'KOTA PADANG',
            ),
            69 =>
            array(
                'id' => '1372',
                'province_id' => '13',
                'name' => 'KOTA SOLOK',
            ),
            70 =>
            array(
                'id' => '1373',
                'province_id' => '13',
                'name' => 'KOTA SAWAHLUNTO',
            ),
            71 =>
            array(
                'id' => '1374',
                'province_id' => '13',
                'name' => 'KOTA PADANG PANJANG',
            ),
            72 =>
            array(
                'id' => '1375',
                'province_id' => '13',
                'name' => 'KOTA BUKITTINGGI',
            ),
            73 =>
            array(
                'id' => '1376',
                'province_id' => '13',
                'name' => 'KOTA PAYAKUMBUH',
            ),
            74 =>
            array(
                'id' => '1377',
                'province_id' => '13',
                'name' => 'KOTA PARIAMAN',
            ),
            75 =>
            array(
                'id' => '1401',
                'province_id' => '14',
                'name' => 'KAB. KAMPAR',
            ),
            76 =>
            array(
                'id' => '1402',
                'province_id' => '14',
                'name' => 'KAB. INDRAGIRI HULU',
            ),
            77 =>
            array(
                'id' => '1403',
                'province_id' => '14',
                'name' => 'KAB. BENGKALIS',
            ),
            78 =>
            array(
                'id' => '1404',
                'province_id' => '14',
                'name' => 'KAB. INDRAGIRI HILIR',
            ),
            79 =>
            array(
                'id' => '1405',
                'province_id' => '14',
                'name' => 'KAB. PELALAWAN',
            ),
            80 =>
            array(
                'id' => '1406',
                'province_id' => '14',
                'name' => 'KAB. ROKAN HULU',
            ),
            81 =>
            array(
                'id' => '1407',
                'province_id' => '14',
                'name' => 'KAB. ROKAN HILIR',
            ),
            82 =>
            array(
                'id' => '1408',
                'province_id' => '14',
                'name' => 'KAB. SIAK',
            ),
            83 =>
            array(
                'id' => '1409',
                'province_id' => '14',
                'name' => 'KAB. KUANTAN SINGINGI',
            ),
            84 =>
            array(
                'id' => '1410',
                'province_id' => '14',
                'name' => 'KAB. KEPULAUAN MERANTI',
            ),
            85 =>
            array(
                'id' => '1471',
                'province_id' => '14',
                'name' => 'KOTA PEKANBARU',
            ),
            86 =>
            array(
                'id' => '1472',
                'province_id' => '14',
                'name' => 'KOTA DUMAI',
            ),
            87 =>
            array(
                'id' => '1501',
                'province_id' => '15',
                'name' => 'KAB. KERINCI',
            ),
            88 =>
            array(
                'id' => '1502',
                'province_id' => '15',
                'name' => 'KAB. MERANGIN',
            ),
            89 =>
            array(
                'id' => '1503',
                'province_id' => '15',
                'name' => 'KAB. SAROLANGUN',
            ),
            90 =>
            array(
                'id' => '1504',
                'province_id' => '15',
                'name' => 'KAB. BATANGHARI',
            ),
            91 =>
            array(
                'id' => '1505',
                'province_id' => '15',
                'name' => 'KAB. MUARO JAMBI',
            ),
            92 =>
            array(
                'id' => '1506',
                'province_id' => '15',
                'name' => 'KAB. TANJUNG JABUNG BARAT',
            ),
            93 =>
            array(
                'id' => '1507',
                'province_id' => '15',
                'name' => 'KAB. TANJUNG JABUNG TIMUR',
            ),
            94 =>
            array(
                'id' => '1508',
                'province_id' => '15',
                'name' => 'KAB. BUNGO',
            ),
            95 =>
            array(
                'id' => '1509',
                'province_id' => '15',
                'name' => 'KAB. TEBO',
            ),
            96 =>
            array(
                'id' => '1571',
                'province_id' => '15',
                'name' => 'KOTA JAMBI',
            ),
            97 =>
            array(
                'id' => '1572',
                'province_id' => '15',
                'name' => 'KOTA SUNGAI PENUH',
            ),
            98 =>
            array(
                'id' => '1601',
                'province_id' => '16',
                'name' => 'KAB. OGAN KOMERING ULU',
            ),
            99 =>
            array(
                'id' => '1602',
                'province_id' => '16',
                'name' => 'KAB. OGAN KOMERING ILIR',
            ),
            100 =>
            array(
                'id' => '1603',
                'province_id' => '16',
                'name' => 'KAB. MUARA ENIM',
            ),
            101 =>
            array(
                'id' => '1604',
                'province_id' => '16',
                'name' => 'KAB. LAHAT',
            ),
            102 =>
            array(
                'id' => '1605',
                'province_id' => '16',
                'name' => 'KAB. MUSI RAWAS',
            ),
            103 =>
            array(
                'id' => '1606',
                'province_id' => '16',
                'name' => 'KAB. MUSI BANYUASIN',
            ),
            104 =>
            array(
                'id' => '1607',
                'province_id' => '16',
                'name' => 'KAB. BANYUASIN',
            ),
            105 =>
            array(
                'id' => '1608',
                'province_id' => '16',
                'name' => 'KAB. OGAN KOMERING ULU TIMUR',
            ),
            106 =>
            array(
                'id' => '1609',
                'province_id' => '16',
                'name' => 'KAB. OGAN KOMERING ULU SELATAN',
            ),
            107 =>
            array(
                'id' => '1610',
                'province_id' => '16',
                'name' => 'KAB. OGAN ILIR',
            ),
            108 =>
            array(
                'id' => '1611',
                'province_id' => '16',
                'name' => 'KAB. EMPAT LAWANG',
            ),
            109 =>
            array(
                'id' => '1612',
                'province_id' => '16',
                'name' => 'KAB. PENUKAL ABAB LEMATANG ILIR',
            ),
            110 =>
            array(
                'id' => '1613',
                'province_id' => '16',
                'name' => 'KAB. MUSI RAWAS UTARA',
            ),
            111 =>
            array(
                'id' => '1671',
                'province_id' => '16',
                'name' => 'KOTA PALEMBANG',
            ),
            112 =>
            array(
                'id' => '1672',
                'province_id' => '16',
                'name' => 'KOTA PAGAR ALAM',
            ),
            113 =>
            array(
                'id' => '1673',
                'province_id' => '16',
                'name' => 'KOTA LUBUK LINGGAU',
            ),
            114 =>
            array(
                'id' => '1674',
                'province_id' => '16',
                'name' => 'KOTA PRABUMULIH',
            ),
            115 =>
            array(
                'id' => '1701',
                'province_id' => '17',
                'name' => 'KAB. BENGKULU SELATAN',
            ),
            116 =>
            array(
                'id' => '1702',
                'province_id' => '17',
                'name' => 'KAB. REJANG LEBONG',
            ),
            117 =>
            array(
                'id' => '1703',
                'province_id' => '17',
                'name' => 'KAB. BENGKULU UTARA',
            ),
            118 =>
            array(
                'id' => '1704',
                'province_id' => '17',
                'name' => 'KAB. KAUR',
            ),
            119 =>
            array(
                'id' => '1705',
                'province_id' => '17',
                'name' => 'KAB. SELUMA',
            ),
            120 =>
            array(
                'id' => '1706',
                'province_id' => '17',
                'name' => 'KAB. MUKO MUKO',
            ),
            121 =>
            array(
                'id' => '1707',
                'province_id' => '17',
                'name' => 'KAB. LEBONG',
            ),
            122 =>
            array(
                'id' => '1708',
                'province_id' => '17',
                'name' => 'KAB. KEPAHIANG',
            ),
            123 =>
            array(
                'id' => '1709',
                'province_id' => '17',
                'name' => 'KAB. BENGKULU TENGAH',
            ),
            124 =>
            array(
                'id' => '1771',
                'province_id' => '17',
                'name' => 'KOTA BENGKULU',
            ),
            125 =>
            array(
                'id' => '1801',
                'province_id' => '18',
                'name' => 'KAB. LAMPUNG SELATAN',
            ),
            126 =>
            array(
                'id' => '1802',
                'province_id' => '18',
                'name' => 'KAB. LAMPUNG TENGAH',
            ),
            127 =>
            array(
                'id' => '1803',
                'province_id' => '18',
                'name' => 'KAB. LAMPUNG UTARA',
            ),
            128 =>
            array(
                'id' => '1804',
                'province_id' => '18',
                'name' => 'KAB. LAMPUNG BARAT',
            ),
            129 =>
            array(
                'id' => '1805',
                'province_id' => '18',
                'name' => 'KAB. TULANG BAWANG',
            ),
            130 =>
            array(
                'id' => '1806',
                'province_id' => '18',
                'name' => 'KAB. TANGGAMUS',
            ),
            131 =>
            array(
                'id' => '1807',
                'province_id' => '18',
                'name' => 'KAB. LAMPUNG TIMUR',
            ),
            132 =>
            array(
                'id' => '1808',
                'province_id' => '18',
                'name' => 'KAB. WAY KANAN',
            ),
            133 =>
            array(
                'id' => '1809',
                'province_id' => '18',
                'name' => 'KAB. PESAWARAN',
            ),
            134 =>
            array(
                'id' => '1810',
                'province_id' => '18',
                'name' => 'KAB. PRINGSEWU',
            ),
            135 =>
            array(
                'id' => '1811',
                'province_id' => '18',
                'name' => 'KAB. MESUJI',
            ),
            136 =>
            array(
                'id' => '1812',
                'province_id' => '18',
                'name' => 'KAB. TULANG BAWANG BARAT',
            ),
            137 =>
            array(
                'id' => '1813',
                'province_id' => '18',
                'name' => 'KAB. PESISIR BARAT',
            ),
            138 =>
            array(
                'id' => '1871',
                'province_id' => '18',
                'name' => 'KOTA BANDAR LAMPUNG',
            ),
            139 =>
            array(
                'id' => '1872',
                'province_id' => '18',
                'name' => 'KOTA METRO',
            ),
            140 =>
            array(
                'id' => '1901',
                'province_id' => '19',
                'name' => 'KAB. BANGKA',
            ),
            141 =>
            array(
                'id' => '1902',
                'province_id' => '19',
                'name' => 'KAB. BELITUNG',
            ),
            142 =>
            array(
                'id' => '1903',
                'province_id' => '19',
                'name' => 'KAB. BANGKA SELATAN',
            ),
            143 =>
            array(
                'id' => '1904',
                'province_id' => '19',
                'name' => 'KAB. BANGKA TENGAH',
            ),
            144 =>
            array(
                'id' => '1905',
                'province_id' => '19',
                'name' => 'KAB. BANGKA BARAT',
            ),
            145 =>
            array(
                'id' => '1906',
                'province_id' => '19',
                'name' => 'KAB. BELITUNG TIMUR',
            ),
            146 =>
            array(
                'id' => '1971',
                'province_id' => '19',
                'name' => 'KOTA PANGKAL PINANG',
            ),
            147 =>
            array(
                'id' => '2101',
                'province_id' => '21',
                'name' => 'KAB. BINTAN',
            ),
            148 =>
            array(
                'id' => '2102',
                'province_id' => '21',
                'name' => 'KAB. KARIMUN',
            ),
            149 =>
            array(
                'id' => '2103',
                'province_id' => '21',
                'name' => 'KAB. NATUNA',
            ),
            150 =>
            array(
                'id' => '2104',
                'province_id' => '21',
                'name' => 'KAB. LINGGA',
            ),
            151 =>
            array(
                'id' => '2105',
                'province_id' => '21',
                'name' => 'KAB. KEPULAUAN ANAMBAS',
            ),
            152 =>
            array(
                'id' => '2171',
                'province_id' => '21',
                'name' => 'KOTA BATAM',
            ),
            153 =>
            array(
                'id' => '2172',
                'province_id' => '21',
                'name' => 'KOTA TANJUNG PINANG',
            ),
            154 =>
            array(
                'id' => '3101',
                'province_id' => '31',
                'name' => 'KAB. ADM. KEP. SERIBU',
            ),
            155 =>
            array(
                'id' => '3171',
                'province_id' => '31',
                'name' => 'KOTA ADM. JAKARTA PUSAT',
            ),
            156 =>
            array(
                'id' => '3172',
                'province_id' => '31',
                'name' => 'KOTA ADM. JAKARTA UTARA',
            ),
            157 =>
            array(
                'id' => '3173',
                'province_id' => '31',
                'name' => 'KOTA ADM. JAKARTA BARAT',
            ),
            158 =>
            array(
                'id' => '3174',
                'province_id' => '31',
                'name' => 'KOTA ADM. JAKARTA SELATAN',
            ),
            159 =>
            array(
                'id' => '3175',
                'province_id' => '31',
                'name' => 'KOTA ADM. JAKARTA TIMUR',
            ),
            160 =>
            array(
                'id' => '3201',
                'province_id' => '32',
                'name' => 'KAB. BOGOR',
            ),
            161 =>
            array(
                'id' => '3202',
                'province_id' => '32',
                'name' => 'KAB. SUKABUMI',
            ),
            162 =>
            array(
                'id' => '3203',
                'province_id' => '32',
                'name' => 'KAB. CIANJUR',
            ),
            163 =>
            array(
                'id' => '3204',
                'province_id' => '32',
                'name' => 'KAB. BANDUNG',
            ),
            164 =>
            array(
                'id' => '3205',
                'province_id' => '32',
                'name' => 'KAB. GARUT',
            ),
            165 =>
            array(
                'id' => '3206',
                'province_id' => '32',
                'name' => 'KAB. TASIKMALAYA',
            ),
            166 =>
            array(
                'id' => '3207',
                'province_id' => '32',
                'name' => 'KAB. CIAMIS',
            ),
            167 =>
            array(
                'id' => '3208',
                'province_id' => '32',
                'name' => 'KAB. KUNINGAN',
            ),
            168 =>
            array(
                'id' => '3209',
                'province_id' => '32',
                'name' => 'KAB. CIREBON',
            ),
            169 =>
            array(
                'id' => '3210',
                'province_id' => '32',
                'name' => 'KAB. MAJALENGKA',
            ),
            170 =>
            array(
                'id' => '3211',
                'province_id' => '32',
                'name' => 'KAB. SUMEDANG',
            ),
            171 =>
            array(
                'id' => '3212',
                'province_id' => '32',
                'name' => 'KAB. INDRAMAYU',
            ),
            172 =>
            array(
                'id' => '3213',
                'province_id' => '32',
                'name' => 'KAB. SUBANG',
            ),
            173 =>
            array(
                'id' => '3214',
                'province_id' => '32',
                'name' => 'KAB. PURWAKARTA',
            ),
            174 =>
            array(
                'id' => '3215',
                'province_id' => '32',
                'name' => 'KAB. KARAWANG',
            ),
            175 =>
            array(
                'id' => '3216',
                'province_id' => '32',
                'name' => 'KAB. BEKASI',
            ),
            176 =>
            array(
                'id' => '3217',
                'province_id' => '32',
                'name' => 'KAB. BANDUNG BARAT',
            ),
            177 =>
            array(
                'id' => '3218',
                'province_id' => '32',
                'name' => 'KAB. PANGANDARAN',
            ),
            178 =>
            array(
                'id' => '3271',
                'province_id' => '32',
                'name' => 'KOTA BOGOR',
            ),
            179 =>
            array(
                'id' => '3272',
                'province_id' => '32',
                'name' => 'KOTA SUKABUMI',
            ),
            180 =>
            array(
                'id' => '3273',
                'province_id' => '32',
                'name' => 'KOTA BANDUNG',
            ),
            181 =>
            array(
                'id' => '3274',
                'province_id' => '32',
                'name' => 'KOTA CIREBON',
            ),
            182 =>
            array(
                'id' => '3275',
                'province_id' => '32',
                'name' => 'KOTA BEKASI',
            ),
            183 =>
            array(
                'id' => '3276',
                'province_id' => '32',
                'name' => 'KOTA DEPOK',
            ),
            184 =>
            array(
                'id' => '3277',
                'province_id' => '32',
                'name' => 'KOTA CIMAHI',
            ),
            185 =>
            array(
                'id' => '3278',
                'province_id' => '32',
                'name' => 'KOTA TASIKMALAYA',
            ),
            186 =>
            array(
                'id' => '3279',
                'province_id' => '32',
                'name' => 'KOTA BANJAR',
            ),
            187 =>
            array(
                'id' => '3301',
                'province_id' => '33',
                'name' => 'KAB. CILACAP',
            ),
            188 =>
            array(
                'id' => '3302',
                'province_id' => '33',
                'name' => 'KAB. BANYUMAS',
            ),
            189 =>
            array(
                'id' => '3303',
                'province_id' => '33',
                'name' => 'KAB. PURBALINGGA',
            ),
            190 =>
            array(
                'id' => '3304',
                'province_id' => '33',
                'name' => 'KAB. BANJARNEGARA',
            ),
            191 =>
            array(
                'id' => '3305',
                'province_id' => '33',
                'name' => 'KAB. KEBUMEN',
            ),
            192 =>
            array(
                'id' => '3306',
                'province_id' => '33',
                'name' => 'KAB. PURWOREJO',
            ),
            193 =>
            array(
                'id' => '3307',
                'province_id' => '33',
                'name' => 'KAB. WONOSOBO',
            ),
            194 =>
            array(
                'id' => '3308',
                'province_id' => '33',
                'name' => 'KAB. MAGELANG',
            ),
            195 =>
            array(
                'id' => '3309',
                'province_id' => '33',
                'name' => 'KAB. BOYOLALI',
            ),
            196 =>
            array(
                'id' => '3310',
                'province_id' => '33',
                'name' => 'KAB. KLATEN',
            ),
            197 =>
            array(
                'id' => '3311',
                'province_id' => '33',
                'name' => 'KAB. SUKOHARJO',
            ),
            198 =>
            array(
                'id' => '3312',
                'province_id' => '33',
                'name' => 'KAB. WONOGIRI',
            ),
            199 =>
            array(
                'id' => '3313',
                'province_id' => '33',
                'name' => 'KAB. KARANGANYAR',
            ),
            200 =>
            array(
                'id' => '3314',
                'province_id' => '33',
                'name' => 'KAB. SRAGEN',
            ),
            201 =>
            array(
                'id' => '3315',
                'province_id' => '33',
                'name' => 'KAB. GROBOGAN',
            ),
            202 =>
            array(
                'id' => '3316',
                'province_id' => '33',
                'name' => 'KAB. BLORA',
            ),
            203 =>
            array(
                'id' => '3317',
                'province_id' => '33',
                'name' => 'KAB. REMBANG',
            ),
            204 =>
            array(
                'id' => '3318',
                'province_id' => '33',
                'name' => 'KAB. PATI',
            ),
            205 =>
            array(
                'id' => '3319',
                'province_id' => '33',
                'name' => 'KAB. KUDUS',
            ),
            206 =>
            array(
                'id' => '3320',
                'province_id' => '33',
                'name' => 'KAB. JEPARA',
            ),
            207 =>
            array(
                'id' => '3321',
                'province_id' => '33',
                'name' => 'KAB. DEMAK',
            ),
            208 =>
            array(
                'id' => '3322',
                'province_id' => '33',
                'name' => 'KAB. SEMARANG',
            ),
            209 =>
            array(
                'id' => '3323',
                'province_id' => '33',
                'name' => 'KAB. TEMANGGUNG',
            ),
            210 =>
            array(
                'id' => '3324',
                'province_id' => '33',
                'name' => 'KAB. KENDAL',
            ),
            211 =>
            array(
                'id' => '3325',
                'province_id' => '33',
                'name' => 'KAB. BATANG',
            ),
            212 =>
            array(
                'id' => '3326',
                'province_id' => '33',
                'name' => 'KAB. PEKALONGAN',
            ),
            213 =>
            array(
                'id' => '3327',
                'province_id' => '33',
                'name' => 'KAB. PEMALANG',
            ),
            214 =>
            array(
                'id' => '3328',
                'province_id' => '33',
                'name' => 'KAB. TEGAL',
            ),
            215 =>
            array(
                'id' => '3329',
                'province_id' => '33',
                'name' => 'KAB. BREBES',
            ),
            216 =>
            array(
                'id' => '3371',
                'province_id' => '33',
                'name' => 'KOTA MAGELANG',
            ),
            217 =>
            array(
                'id' => '3372',
                'province_id' => '33',
                'name' => 'KOTA SURAKARTA',
            ),
            218 =>
            array(
                'id' => '3373',
                'province_id' => '33',
                'name' => 'KOTA SALATIGA',
            ),
            219 =>
            array(
                'id' => '3374',
                'province_id' => '33',
                'name' => 'KOTA SEMARANG',
            ),
            220 =>
            array(
                'id' => '3375',
                'province_id' => '33',
                'name' => 'KOTA PEKALONGAN',
            ),
            221 =>
            array(
                'id' => '3376',
                'province_id' => '33',
                'name' => 'KOTA TEGAL',
            ),
            222 =>
            array(
                'id' => '3401',
                'province_id' => '34',
                'name' => 'KAB. KULON PROGO',
            ),
            223 =>
            array(
                'id' => '3402',
                'province_id' => '34',
                'name' => 'KAB. BANTUL',
            ),
            224 =>
            array(
                'id' => '3403',
                'province_id' => '34',
                'name' => 'KAB. GUNUNGKIDUL',
            ),
            225 =>
            array(
                'id' => '3404',
                'province_id' => '34',
                'name' => 'KAB. SLEMAN',
            ),
            226 =>
            array(
                'id' => '3471',
                'province_id' => '34',
                'name' => 'KOTA YOGYAKARTA',
            ),
            227 =>
            array(
                'id' => '3501',
                'province_id' => '35',
                'name' => 'KAB. PACITAN',
            ),
            228 =>
            array(
                'id' => '3502',
                'province_id' => '35',
                'name' => 'KAB. PONOROGO',
            ),
            229 =>
            array(
                'id' => '3503',
                'province_id' => '35',
                'name' => 'KAB. TRENGGALEK',
            ),
            230 =>
            array(
                'id' => '3504',
                'province_id' => '35',
                'name' => 'KAB. TULUNGAGUNG',
            ),
            231 =>
            array(
                'id' => '3505',
                'province_id' => '35',
                'name' => 'KAB. BLITAR',
            ),
            232 =>
            array(
                'id' => '3506',
                'province_id' => '35',
                'name' => 'KAB. KEDIRI',
            ),
            233 =>
            array(
                'id' => '3507',
                'province_id' => '35',
                'name' => 'KAB. MALANG',
            ),
            234 =>
            array(
                'id' => '3508',
                'province_id' => '35',
                'name' => 'KAB. LUMAJANG',
            ),
            235 =>
            array(
                'id' => '3509',
                'province_id' => '35',
                'name' => 'KAB. JEMBER',
            ),
            236 =>
            array(
                'id' => '3510',
                'province_id' => '35',
                'name' => 'KAB. BANYUWANGI',
            ),
            237 =>
            array(
                'id' => '3511',
                'province_id' => '35',
                'name' => 'KAB. BONDOWOSO',
            ),
            238 =>
            array(
                'id' => '3512',
                'province_id' => '35',
                'name' => 'KAB. SITUBONDO',
            ),
            239 =>
            array(
                'id' => '3513',
                'province_id' => '35',
                'name' => 'KAB. PROBOLINGGO',
            ),
            240 =>
            array(
                'id' => '3514',
                'province_id' => '35',
                'name' => 'KAB. PASURUAN',
            ),
            241 =>
            array(
                'id' => '3515',
                'province_id' => '35',
                'name' => 'KAB. SIDOARJO',
            ),
            242 =>
            array(
                'id' => '3516',
                'province_id' => '35',
                'name' => 'KAB. MOJOKERTO',
            ),
            243 =>
            array(
                'id' => '3517',
                'province_id' => '35',
                'name' => 'KAB. JOMBANG',
            ),
            244 =>
            array(
                'id' => '3518',
                'province_id' => '35',
                'name' => 'KAB. NGANJUK',
            ),
            245 =>
            array(
                'id' => '3519',
                'province_id' => '35',
                'name' => 'KAB. MADIUN',
            ),
            246 =>
            array(
                'id' => '3520',
                'province_id' => '35',
                'name' => 'KAB. MAGETAN',
            ),
            247 =>
            array(
                'id' => '3521',
                'province_id' => '35',
                'name' => 'KAB. NGAWI',
            ),
            248 =>
            array(
                'id' => '3522',
                'province_id' => '35',
                'name' => 'KAB. BOJONEGORO',
            ),
            249 =>
            array(
                'id' => '3523',
                'province_id' => '35',
                'name' => 'KAB. TUBAN',
            ),
            250 =>
            array(
                'id' => '3524',
                'province_id' => '35',
                'name' => 'KAB. LAMONGAN',
            ),
            251 =>
            array(
                'id' => '3525',
                'province_id' => '35',
                'name' => 'KAB. GRESIK',
            ),
            252 =>
            array(
                'id' => '3526',
                'province_id' => '35',
                'name' => 'KAB. BANGKALAN',
            ),
            253 =>
            array(
                'id' => '3527',
                'province_id' => '35',
                'name' => 'KAB. SAMPANG',
            ),
            254 =>
            array(
                'id' => '3528',
                'province_id' => '35',
                'name' => 'KAB. PAMEKASAN',
            ),
            255 =>
            array(
                'id' => '3529',
                'province_id' => '35',
                'name' => 'KAB. SUMENEP',
            ),
            256 =>
            array(
                'id' => '3571',
                'province_id' => '35',
                'name' => 'KOTA KEDIRI',
            ),
            257 =>
            array(
                'id' => '3572',
                'province_id' => '35',
                'name' => 'KOTA BLITAR',
            ),
            258 =>
            array(
                'id' => '3573',
                'province_id' => '35',
                'name' => 'KOTA MALANG',
            ),
            259 =>
            array(
                'id' => '3574',
                'province_id' => '35',
                'name' => 'KOTA PROBOLINGGO',
            ),
            260 =>
            array(
                'id' => '3575',
                'province_id' => '35',
                'name' => 'KOTA PASURUAN',
            ),
            261 =>
            array(
                'id' => '3576',
                'province_id' => '35',
                'name' => 'KOTA MOJOKERTO',
            ),
            262 =>
            array(
                'id' => '3577',
                'province_id' => '35',
                'name' => 'KOTA MADIUN',
            ),
            263 =>
            array(
                'id' => '3578',
                'province_id' => '35',
                'name' => 'KOTA SURABAYA',
            ),
            264 =>
            array(
                'id' => '3579',
                'province_id' => '35',
                'name' => 'KOTA BATU',
            ),
            265 =>
            array(
                'id' => '3601',
                'province_id' => '36',
                'name' => 'KAB. PANDEGLANG',
            ),
            266 =>
            array(
                'id' => '3602',
                'province_id' => '36',
                'name' => 'KAB. LEBAK',
            ),
            267 =>
            array(
                'id' => '3603',
                'province_id' => '36',
                'name' => 'KAB. TANGERANG',
            ),
            268 =>
            array(
                'id' => '3604',
                'province_id' => '36',
                'name' => 'KAB. SERANG',
            ),
            269 =>
            array(
                'id' => '3671',
                'province_id' => '36',
                'name' => 'KOTA TANGERANG',
            ),
            270 =>
            array(
                'id' => '3672',
                'province_id' => '36',
                'name' => 'KOTA CILEGON',
            ),
            271 =>
            array(
                'id' => '3673',
                'province_id' => '36',
                'name' => 'KOTA SERANG',
            ),
            272 =>
            array(
                'id' => '3674',
                'province_id' => '36',
                'name' => 'KOTA TANGERANG SELATAN',
            ),
            273 =>
            array(
                'id' => '5101',
                'province_id' => '51',
                'name' => 'KAB. JEMBRANA',
            ),
            274 =>
            array(
                'id' => '5102',
                'province_id' => '51',
                'name' => 'KAB. TABANAN',
            ),
            275 =>
            array(
                'id' => '5103',
                'province_id' => '51',
                'name' => 'KAB. BADUNG',
            ),
            276 =>
            array(
                'id' => '5104',
                'province_id' => '51',
                'name' => 'KAB. GIANYAR',
            ),
            277 =>
            array(
                'id' => '5105',
                'province_id' => '51',
                'name' => 'KAB. KLUNGKUNG',
            ),
            278 =>
            array(
                'id' => '5106',
                'province_id' => '51',
                'name' => 'KAB. BANGLI',
            ),
            279 =>
            array(
                'id' => '5107',
                'province_id' => '51',
                'name' => 'KAB. KARANGASEM',
            ),
            280 =>
            array(
                'id' => '5108',
                'province_id' => '51',
                'name' => 'KAB. BULELENG',
            ),
            281 =>
            array(
                'id' => '5171',
                'province_id' => '51',
                'name' => 'KOTA DENPASAR',
            ),
            282 =>
            array(
                'id' => '5201',
                'province_id' => '52',
                'name' => 'KAB. LOMBOK BARAT',
            ),
            283 =>
            array(
                'id' => '5202',
                'province_id' => '52',
                'name' => 'KAB. LOMBOK TENGAH',
            ),
            284 =>
            array(
                'id' => '5203',
                'province_id' => '52',
                'name' => 'KAB. LOMBOK TIMUR',
            ),
            285 =>
            array(
                'id' => '5204',
                'province_id' => '52',
                'name' => 'KAB. SUMBAWA',
            ),
            286 =>
            array(
                'id' => '5205',
                'province_id' => '52',
                'name' => 'KAB. DOMPU',
            ),
            287 =>
            array(
                'id' => '5206',
                'province_id' => '52',
                'name' => 'KAB. BIMA',
            ),
            288 =>
            array(
                'id' => '5207',
                'province_id' => '52',
                'name' => 'KAB. SUMBAWA BARAT',
            ),
            289 =>
            array(
                'id' => '5208',
                'province_id' => '52',
                'name' => 'KAB. LOMBOK UTARA',
            ),
            290 =>
            array(
                'id' => '5271',
                'province_id' => '52',
                'name' => 'KOTA MATARAM',
            ),
            291 =>
            array(
                'id' => '5272',
                'province_id' => '52',
                'name' => 'KOTA BIMA',
            ),
            292 =>
            array(
                'id' => '5301',
                'province_id' => '53',
                'name' => 'KAB. KUPANG',
            ),
            293 =>
            array(
                'id' => '5302',
                'province_id' => '53',
                'name' => 'KAB TIMOR TENGAH SELATAN',
            ),
            294 =>
            array(
                'id' => '5303',
                'province_id' => '53',
                'name' => 'KAB. TIMOR TENGAH UTARA',
            ),
            295 =>
            array(
                'id' => '5304',
                'province_id' => '53',
                'name' => 'KAB. BELU',
            ),
            296 =>
            array(
                'id' => '5305',
                'province_id' => '53',
                'name' => 'KAB. ALOR',
            ),
            297 =>
            array(
                'id' => '5306',
                'province_id' => '53',
                'name' => 'KAB. FLORES TIMUR',
            ),
            298 =>
            array(
                'id' => '5307',
                'province_id' => '53',
                'name' => 'KAB. SIKKA',
            ),
            299 =>
            array(
                'id' => '5308',
                'province_id' => '53',
                'name' => 'KAB. ENDE',
            ),
            300 =>
            array(
                'id' => '5309',
                'province_id' => '53',
                'name' => 'KAB. NGADA',
            ),
            301 =>
            array(
                'id' => '5310',
                'province_id' => '53',
                'name' => 'KAB. MANGGARAI',
            ),
            302 =>
            array(
                'id' => '5311',
                'province_id' => '53',
                'name' => 'KAB. SUMBA TIMUR',
            ),
            303 =>
            array(
                'id' => '5312',
                'province_id' => '53',
                'name' => 'KAB. SUMBA BARAT',
            ),
            304 =>
            array(
                'id' => '5313',
                'province_id' => '53',
                'name' => 'KAB. LEMBATA',
            ),
            305 =>
            array(
                'id' => '5314',
                'province_id' => '53',
                'name' => 'KAB. ROTE NDAO',
            ),
            306 =>
            array(
                'id' => '5315',
                'province_id' => '53',
                'name' => 'KAB. MANGGARAI BARAT',
            ),
            307 =>
            array(
                'id' => '5316',
                'province_id' => '53',
                'name' => 'KAB. NAGEKEO',
            ),
            308 =>
            array(
                'id' => '5317',
                'province_id' => '53',
                'name' => 'KAB. SUMBA TENGAH',
            ),
            309 =>
            array(
                'id' => '5318',
                'province_id' => '53',
                'name' => 'KAB. SUMBA BARAT DAYA',
            ),
            310 =>
            array(
                'id' => '5319',
                'province_id' => '53',
                'name' => 'KAB. MANGGARAI TIMUR',
            ),
            311 =>
            array(
                'id' => '5320',
                'province_id' => '53',
                'name' => 'KAB. SABU RAIJUA',
            ),
            312 =>
            array(
                'id' => '5321',
                'province_id' => '53',
                'name' => 'KAB. MALAKA',
            ),
            313 =>
            array(
                'id' => '5371',
                'province_id' => '53',
                'name' => 'KOTA KUPANG',
            ),
            314 =>
            array(
                'id' => '6101',
                'province_id' => '61',
                'name' => 'KAB. SAMBAS',
            ),
            315 =>
            array(
                'id' => '6102',
                'province_id' => '61',
                'name' => 'KAB. MEMPAWAH',
            ),
            316 =>
            array(
                'id' => '6103',
                'province_id' => '61',
                'name' => 'KAB. SANGGAU',
            ),
            317 =>
            array(
                'id' => '6104',
                'province_id' => '61',
                'name' => 'KAB. KETAPANG',
            ),
            318 =>
            array(
                'id' => '6105',
                'province_id' => '61',
                'name' => 'KAB. SINTANG',
            ),
            319 =>
            array(
                'id' => '6106',
                'province_id' => '61',
                'name' => 'KAB. KAPUAS HULU',
            ),
            320 =>
            array(
                'id' => '6107',
                'province_id' => '61',
                'name' => 'KAB. BENGKAYANG',
            ),
            321 =>
            array(
                'id' => '6108',
                'province_id' => '61',
                'name' => 'KAB. LANDAK',
            ),
            322 =>
            array(
                'id' => '6109',
                'province_id' => '61',
                'name' => 'KAB. SEKADAU',
            ),
            323 =>
            array(
                'id' => '6110',
                'province_id' => '61',
                'name' => 'KAB. MELAWI',
            ),
            324 =>
            array(
                'id' => '6111',
                'province_id' => '61',
                'name' => 'KAB. KAYONG UTARA',
            ),
            325 =>
            array(
                'id' => '6112',
                'province_id' => '61',
                'name' => 'KAB. KUBU RAYA',
            ),
            326 =>
            array(
                'id' => '6171',
                'province_id' => '61',
                'name' => 'KOTA PONTIANAK',
            ),
            327 =>
            array(
                'id' => '6172',
                'province_id' => '61',
                'name' => 'KOTA SINGKAWANG',
            ),
            328 =>
            array(
                'id' => '6201',
                'province_id' => '62',
                'name' => 'KAB. KOTAWARINGIN BARAT',
            ),
            329 =>
            array(
                'id' => '6202',
                'province_id' => '62',
                'name' => 'KAB. KOTAWARINGIN TIMUR',
            ),
            330 =>
            array(
                'id' => '6203',
                'province_id' => '62',
                'name' => 'KAB. KAPUAS',
            ),
            331 =>
            array(
                'id' => '6204',
                'province_id' => '62',
                'name' => 'KAB. BARITO SELATAN',
            ),
            332 =>
            array(
                'id' => '6205',
                'province_id' => '62',
                'name' => 'KAB. BARITO UTARA',
            ),
            333 =>
            array(
                'id' => '6206',
                'province_id' => '62',
                'name' => 'KAB. KATINGAN',
            ),
            334 =>
            array(
                'id' => '6207',
                'province_id' => '62',
                'name' => 'KAB. SERUYAN',
            ),
            335 =>
            array(
                'id' => '6208',
                'province_id' => '62',
                'name' => 'KAB. SUKAMARA',
            ),
            336 =>
            array(
                'id' => '6209',
                'province_id' => '62',
                'name' => 'KAB. LAMANDAU',
            ),
            337 =>
            array(
                'id' => '6210',
                'province_id' => '62',
                'name' => 'KAB. GUNUNG MAS',
            ),
            338 =>
            array(
                'id' => '6211',
                'province_id' => '62',
                'name' => 'KAB. PULANG PISAU',
            ),
            339 =>
            array(
                'id' => '6212',
                'province_id' => '62',
                'name' => 'KAB. MURUNG RAYA',
            ),
            340 =>
            array(
                'id' => '6213',
                'province_id' => '62',
                'name' => 'KAB. BARITO TIMUR',
            ),
            341 =>
            array(
                'id' => '6271',
                'province_id' => '62',
                'name' => 'KOTA PALANGKARAYA',
            ),
            342 =>
            array(
                'id' => '6301',
                'province_id' => '63',
                'name' => 'KAB. TANAH LAUT',
            ),
            343 =>
            array(
                'id' => '6302',
                'province_id' => '63',
                'name' => 'KAB. KOTABARU',
            ),
            344 =>
            array(
                'id' => '6303',
                'province_id' => '63',
                'name' => 'KAB. BANJAR',
            ),
            345 =>
            array(
                'id' => '6304',
                'province_id' => '63',
                'name' => 'KAB. BARITO KUALA',
            ),
            346 =>
            array(
                'id' => '6305',
                'province_id' => '63',
                'name' => 'KAB. TAPIN',
            ),
            347 =>
            array(
                'id' => '6306',
                'province_id' => '63',
                'name' => 'KAB. HULU SUNGAI SELATAN',
            ),
            348 =>
            array(
                'id' => '6307',
                'province_id' => '63',
                'name' => 'KAB. HULU SUNGAI TENGAH',
            ),
            349 =>
            array(
                'id' => '6308',
                'province_id' => '63',
                'name' => 'KAB. HULU SUNGAI UTARA',
            ),
            350 =>
            array(
                'id' => '6309',
                'province_id' => '63',
                'name' => 'KAB. TABALONG',
            ),
            351 =>
            array(
                'id' => '6310',
                'province_id' => '63',
                'name' => 'KAB. TANAH BUMBU',
            ),
            352 =>
            array(
                'id' => '6311',
                'province_id' => '63',
                'name' => 'KAB. BALANGAN',
            ),
            353 =>
            array(
                'id' => '6371',
                'province_id' => '63',
                'name' => 'KOTA BANJARMASIN',
            ),
            354 =>
            array(
                'id' => '6372',
                'province_id' => '63',
                'name' => 'KOTA BANJARBARU',
            ),
            355 =>
            array(
                'id' => '6401',
                'province_id' => '64',
                'name' => 'KAB. PASER',
            ),
            356 =>
            array(
                'id' => '6402',
                'province_id' => '64',
                'name' => 'KAB. KUTAI KARTANEGARA',
            ),
            357 =>
            array(
                'id' => '6403',
                'province_id' => '64',
                'name' => 'KAB. BERAU',
            ),
            358 =>
            array(
                'id' => '6407',
                'province_id' => '64',
                'name' => 'KAB. KUTAI BARAT',
            ),
            359 =>
            array(
                'id' => '6408',
                'province_id' => '64',
                'name' => 'KAB. KUTAI TIMUR',
            ),
            360 =>
            array(
                'id' => '6409',
                'province_id' => '64',
                'name' => 'KAB. PENAJAM PASER UTARA',
            ),
            361 =>
            array(
                'id' => '6411',
                'province_id' => '64',
                'name' => 'KAB. MAHAKAM ULU',
            ),
            362 =>
            array(
                'id' => '6471',
                'province_id' => '64',
                'name' => 'KOTA BALIKPAPAN',
            ),
            363 =>
            array(
                'id' => '6472',
                'province_id' => '64',
                'name' => 'KOTA SAMARINDA',
            ),
            364 =>
            array(
                'id' => '6474',
                'province_id' => '64',
                'name' => 'KOTA BONTANG',
            ),
            365 =>
            array(
                'id' => '6501',
                'province_id' => '65',
                'name' => 'KAB. BULUNGAN',
            ),
            366 =>
            array(
                'id' => '6502',
                'province_id' => '65',
                'name' => 'KAB. MALINAU',
            ),
            367 =>
            array(
                'id' => '6503',
                'province_id' => '65',
                'name' => 'KAB. NUNUKAN',
            ),
            368 =>
            array(
                'id' => '6504',
                'province_id' => '65',
                'name' => 'KAB. TANA TIDUNG',
            ),
            369 =>
            array(
                'id' => '6571',
                'province_id' => '65',
                'name' => 'KOTA TARAKAN',
            ),
            370 =>
            array(
                'id' => '7101',
                'province_id' => '71',
                'name' => 'KAB. BOLAANG MONGONDOW',
            ),
            371 =>
            array(
                'id' => '7102',
                'province_id' => '71',
                'name' => 'KAB. MINAHASA',
            ),
            372 =>
            array(
                'id' => '7103',
                'province_id' => '71',
                'name' => 'KAB. KEPULAUAN SANGIHE',
            ),
            373 =>
            array(
                'id' => '7104',
                'province_id' => '71',
                'name' => 'KAB. KEPULAUAN TALAUD',
            ),
            374 =>
            array(
                'id' => '7105',
                'province_id' => '71',
                'name' => 'KAB. MINAHASA SELATAN',
            ),
            375 =>
            array(
                'id' => '7106',
                'province_id' => '71',
                'name' => 'KAB. MINAHASA UTARA',
            ),
            376 =>
            array(
                'id' => '7107',
                'province_id' => '71',
                'name' => 'KAB. MINAHASA TENGGARA',
            ),
            377 =>
            array(
                'id' => '7108',
                'province_id' => '71',
                'name' => 'KAB. BOLAANG MONGONDOW UTARA',
            ),
            378 =>
            array(
                'id' => '7109',
                'province_id' => '71',
                'name' => 'KAB. KEP. SIAU TAGULANDANG BIARO',
            ),
            379 =>
            array(
                'id' => '7110',
                'province_id' => '71',
                'name' => 'KAB. BOLAANG MONGONDOW TIMUR',
            ),
            380 =>
            array(
                'id' => '7111',
                'province_id' => '71',
                'name' => 'KAB. BOLAANG MONGONDOW SELATAN',
            ),
            381 =>
            array(
                'id' => '7171',
                'province_id' => '71',
                'name' => 'KOTA MANADO',
            ),
            382 =>
            array(
                'id' => '7172',
                'province_id' => '71',
                'name' => 'KOTA BITUNG',
            ),
            383 =>
            array(
                'id' => '7173',
                'province_id' => '71',
                'name' => 'KOTA TOMOHON',
            ),
            384 =>
            array(
                'id' => '7174',
                'province_id' => '71',
                'name' => 'KOTA KOTAMOBAGU',
            ),
            385 =>
            array(
                'id' => '7201',
                'province_id' => '72',
                'name' => 'KAB. BANGGAI',
            ),
            386 =>
            array(
                'id' => '7202',
                'province_id' => '72',
                'name' => 'KAB. POSO',
            ),
            387 =>
            array(
                'id' => '7203',
                'province_id' => '72',
                'name' => 'KAB. DONGGALA',
            ),
            388 =>
            array(
                'id' => '7204',
                'province_id' => '72',
                'name' => 'KAB. TOLI TOLI',
            ),
            389 =>
            array(
                'id' => '7205',
                'province_id' => '72',
                'name' => 'KAB. BUOL',
            ),
            390 =>
            array(
                'id' => '7206',
                'province_id' => '72',
                'name' => 'KAB. MOROWALI',
            ),
            391 =>
            array(
                'id' => '7207',
                'province_id' => '72',
                'name' => 'KAB. BANGGAI KEPULAUAN',
            ),
            392 =>
            array(
                'id' => '7208',
                'province_id' => '72',
                'name' => 'KAB. PARIGI MOUTONG',
            ),
            393 =>
            array(
                'id' => '7209',
                'province_id' => '72',
                'name' => 'KAB. TOJO UNA UNA',
            ),
            394 =>
            array(
                'id' => '7210',
                'province_id' => '72',
                'name' => 'KAB. SIGI',
            ),
            395 =>
            array(
                'id' => '7211',
                'province_id' => '72',
                'name' => 'KAB. BANGGAI LAUT',
            ),
            396 =>
            array(
                'id' => '7212',
                'province_id' => '72',
                'name' => 'KAB. MOROWALI UTARA',
            ),
            397 =>
            array(
                'id' => '7271',
                'province_id' => '72',
                'name' => 'KOTA PALU',
            ),
            398 =>
            array(
                'id' => '7301',
                'province_id' => '73',
                'name' => 'KAB. KEPULAUAN SELAYAR',
            ),
            399 =>
            array(
                'id' => '7302',
                'province_id' => '73',
                'name' => 'KAB. BULUKUMBA',
            ),
            400 =>
            array(
                'id' => '7303',
                'province_id' => '73',
                'name' => 'KAB. BANTAENG',
            ),
            401 =>
            array(
                'id' => '7304',
                'province_id' => '73',
                'name' => 'KAB. JENEPONTO',
            ),
            402 =>
            array(
                'id' => '7305',
                'province_id' => '73',
                'name' => 'KAB. TAKALAR',
            ),
            403 =>
            array(
                'id' => '7306',
                'province_id' => '73',
                'name' => 'KAB. GOWA',
            ),
            404 =>
            array(
                'id' => '7307',
                'province_id' => '73',
                'name' => 'KAB. SINJAI',
            ),
            405 =>
            array(
                'id' => '7308',
                'province_id' => '73',
                'name' => 'KAB. BONE',
            ),
            406 =>
            array(
                'id' => '7309',
                'province_id' => '73',
                'name' => 'KAB. MAROS',
            ),
            407 =>
            array(
                'id' => '7310',
                'province_id' => '73',
                'name' => 'KAB. PANGKAJENE KEPULAUAN',
            ),
            408 =>
            array(
                'id' => '7311',
                'province_id' => '73',
                'name' => 'KAB. BARRU',
            ),
            409 =>
            array(
                'id' => '7312',
                'province_id' => '73',
                'name' => 'KAB. SOPPENG',
            ),
            410 =>
            array(
                'id' => '7313',
                'province_id' => '73',
                'name' => 'KAB. WAJO',
            ),
            411 =>
            array(
                'id' => '7314',
                'province_id' => '73',
                'name' => 'KAB. SIDENRENG RAPPANG',
            ),
            412 =>
            array(
                'id' => '7315',
                'province_id' => '73',
                'name' => 'KAB. PINRANG',
            ),
            413 =>
            array(
                'id' => '7316',
                'province_id' => '73',
                'name' => 'KAB. ENREKANG',
            ),
            414 =>
            array(
                'id' => '7317',
                'province_id' => '73',
                'name' => 'KAB. LUWU',
            ),
            415 =>
            array(
                'id' => '7318',
                'province_id' => '73',
                'name' => 'KAB. TANA TORAJA',
            ),
            416 =>
            array(
                'id' => '7322',
                'province_id' => '73',
                'name' => 'KAB. LUWU UTARA',
            ),
            417 =>
            array(
                'id' => '7324',
                'province_id' => '73',
                'name' => 'KAB. LUWU TIMUR',
            ),
            418 =>
            array(
                'id' => '7326',
                'province_id' => '73',
                'name' => 'KAB. TORAJA UTARA',
            ),
            419 =>
            array(
                'id' => '7371',
                'province_id' => '73',
                'name' => 'KOTA MAKASSAR',
            ),
            420 =>
            array(
                'id' => '7372',
                'province_id' => '73',
                'name' => 'KOTA PARE PARE',
            ),
            421 =>
            array(
                'id' => '7373',
                'province_id' => '73',
                'name' => 'KOTA PALOPO',
            ),
            422 =>
            array(
                'id' => '7401',
                'province_id' => '74',
                'name' => 'KAB. KOLAKA',
            ),
            423 =>
            array(
                'id' => '7402',
                'province_id' => '74',
                'name' => 'KAB. KONAWE',
            ),
            424 =>
            array(
                'id' => '7403',
                'province_id' => '74',
                'name' => 'KAB. MUNA',
            ),
            425 =>
            array(
                'id' => '7404',
                'province_id' => '74',
                'name' => 'KAB. BUTON',
            ),
            426 =>
            array(
                'id' => '7405',
                'province_id' => '74',
                'name' => 'KAB. KONAWE SELATAN',
            ),
            427 =>
            array(
                'id' => '7406',
                'province_id' => '74',
                'name' => 'KAB. BOMBANA',
            ),
            428 =>
            array(
                'id' => '7407',
                'province_id' => '74',
                'name' => 'KAB. WAKATOBI',
            ),
            429 =>
            array(
                'id' => '7408',
                'province_id' => '74',
                'name' => 'KAB. KOLAKA UTARA',
            ),
            430 =>
            array(
                'id' => '7409',
                'province_id' => '74',
                'name' => 'KAB. KONAWE UTARA',
            ),
            431 =>
            array(
                'id' => '7410',
                'province_id' => '74',
                'name' => 'KAB. BUTON UTARA',
            ),
            432 =>
            array(
                'id' => '7411',
                'province_id' => '74',
                'name' => 'KAB. KOLAKA TIMUR',
            ),
            433 =>
            array(
                'id' => '7412',
                'province_id' => '74',
                'name' => 'KAB. KONAWE KEPULAUAN',
            ),
            434 =>
            array(
                'id' => '7413',
                'province_id' => '74',
                'name' => 'KAB. MUNA BARAT',
            ),
            435 =>
            array(
                'id' => '7414',
                'province_id' => '74',
                'name' => 'KAB. BUTON TENGAH',
            ),
            436 =>
            array(
                'id' => '7415',
                'province_id' => '74',
                'name' => 'KAB. BUTON SELATAN',
            ),
            437 =>
            array(
                'id' => '7471',
                'province_id' => '74',
                'name' => 'KOTA KENDARI',
            ),
            438 =>
            array(
                'id' => '7472',
                'province_id' => '74',
                'name' => 'KOTA BAU BAU',
            ),
            439 =>
            array(
                'id' => '7501',
                'province_id' => '75',
                'name' => 'KAB. GORONTALO',
            ),
            440 =>
            array(
                'id' => '7502',
                'province_id' => '75',
                'name' => 'KAB. BOALEMO',
            ),
            441 =>
            array(
                'id' => '7503',
                'province_id' => '75',
                'name' => 'KAB. BONE BOLANGO',
            ),
            442 =>
            array(
                'id' => '7504',
                'province_id' => '75',
                'name' => 'KAB. PAHUWATO',
            ),
            443 =>
            array(
                'id' => '7505',
                'province_id' => '75',
                'name' => 'KAB. GORONTALO UTARA',
            ),
            444 =>
            array(
                'id' => '7571',
                'province_id' => '75',
                'name' => 'KOTA GORONTALO',
            ),
            445 =>
            array(
                'id' => '7601',
                'province_id' => '76',
                'name' => 'KAB. PASANGKAYU',
            ),
            446 =>
            array(
                'id' => '7602',
                'province_id' => '76',
                'name' => 'KAB. MAMUJU',
            ),
            447 =>
            array(
                'id' => '7603',
                'province_id' => '76',
                'name' => 'KAB. MAMASA',
            ),
            448 =>
            array(
                'id' => '7604',
                'province_id' => '76',
                'name' => 'KAB. POLEWALI MANDAR',
            ),
            449 =>
            array(
                'id' => '7605',
                'province_id' => '76',
                'name' => 'KAB. MAJENE',
            ),
            450 =>
            array(
                'id' => '7606',
                'province_id' => '76',
                'name' => 'KAB. MAMUJU TENGAH',
            ),
            451 =>
            array(
                'id' => '8101',
                'province_id' => '81',
                'name' => 'KAB. MALUKU TENGAH',
            ),
            452 =>
            array(
                'id' => '8102',
                'province_id' => '81',
                'name' => 'KAB. MALUKU TENGGARA',
            ),
            453 =>
            array(
                'id' => '8103',
                'province_id' => '81',
                'name' => 'KAB. KEPULAUAN TANIMBAR',
            ),
            454 =>
            array(
                'id' => '8104',
                'province_id' => '81',
                'name' => 'KAB. BURU',
            ),
            455 =>
            array(
                'id' => '8105',
                'province_id' => '81',
                'name' => 'KAB. SERAM BAGIAN TIMUR',
            ),
            456 =>
            array(
                'id' => '8106',
                'province_id' => '81',
                'name' => 'KAB. SERAM BAGIAN BARAT',
            ),
            457 =>
            array(
                'id' => '8107',
                'province_id' => '81',
                'name' => 'KAB. KEPULAUAN ARU',
            ),
            458 =>
            array(
                'id' => '8108',
                'province_id' => '81',
                'name' => 'KAB. MALUKU BARAT DAYA',
            ),
            459 =>
            array(
                'id' => '8109',
                'province_id' => '81',
                'name' => 'KAB. BURU SELATAN',
            ),
            460 =>
            array(
                'id' => '8171',
                'province_id' => '81',
                'name' => 'KOTA AMBON',
            ),
            461 =>
            array(
                'id' => '8172',
                'province_id' => '81',
                'name' => 'KOTA TUAL',
            ),
            462 =>
            array(
                'id' => '8201',
                'province_id' => '82',
                'name' => 'KAB. HALMAHERA BARAT',
            ),
            463 =>
            array(
                'id' => '8202',
                'province_id' => '82',
                'name' => 'KAB. HALMAHERA TENGAH',
            ),
            464 =>
            array(
                'id' => '8203',
                'province_id' => '82',
                'name' => 'KAB. HALMAHERA UTARA',
            ),
            465 =>
            array(
                'id' => '8204',
                'province_id' => '82',
                'name' => 'KAB. HALMAHERA SELATAN',
            ),
            466 =>
            array(
                'id' => '8205',
                'province_id' => '82',
                'name' => 'KAB. KEPULAUAN SULA',
            ),
            467 =>
            array(
                'id' => '8206',
                'province_id' => '82',
                'name' => 'KAB. HALMAHERA TIMUR',
            ),
            468 =>
            array(
                'id' => '8207',
                'province_id' => '82',
                'name' => 'KAB. PULAU MOROTAI',
            ),
            469 =>
            array(
                'id' => '8208',
                'province_id' => '82',
                'name' => 'KAB. PULAU TALIABU',
            ),
            470 =>
            array(
                'id' => '8271',
                'province_id' => '82',
                'name' => 'KOTA TERNATE',
            ),
            471 =>
            array(
                'id' => '8272',
                'province_id' => '82',
                'name' => 'KOTA TIDORE KEPULAUAN',
            ),
            472 =>
            array(
                'id' => '9101',
                'province_id' => '91',
                'name' => 'KAB. MERAUKE',
            ),
            473 =>
            array(
                'id' => '9102',
                'province_id' => '91',
                'name' => 'KAB. JAYAWIJAYA',
            ),
            474 =>
            array(
                'id' => '9103',
                'province_id' => '91',
                'name' => 'KAB. JAYAPURA',
            ),
            475 =>
            array(
                'id' => '9104',
                'province_id' => '91',
                'name' => 'KAB. NABIRE',
            ),
            476 =>
            array(
                'id' => '9105',
                'province_id' => '91',
                'name' => 'KAB. KEPULAUAN YAPEN',
            ),
            477 =>
            array(
                'id' => '9106',
                'province_id' => '91',
                'name' => 'KAB. BIAK NUMFOR',
            ),
            478 =>
            array(
                'id' => '9107',
                'province_id' => '91',
                'name' => 'KAB. PUNCAK JAYA',
            ),
            479 =>
            array(
                'id' => '9108',
                'province_id' => '91',
                'name' => 'KAB. PANIAI',
            ),
            480 =>
            array(
                'id' => '9109',
                'province_id' => '91',
                'name' => 'KAB. MIMIKA',
            ),
            481 =>
            array(
                'id' => '9110',
                'province_id' => '91',
                'name' => 'KAB. SARMI',
            ),
            482 =>
            array(
                'id' => '9111',
                'province_id' => '91',
                'name' => 'KAB. KEEROM',
            ),
            483 =>
            array(
                'id' => '9112',
                'province_id' => '91',
                'name' => 'KAB PEGUNUNGAN BINTANG',
            ),
            484 =>
            array(
                'id' => '9113',
                'province_id' => '91',
                'name' => 'KAB. YAHUKIMO',
            ),
            485 =>
            array(
                'id' => '9114',
                'province_id' => '91',
                'name' => 'KAB. TOLIKARA',
            ),
            486 =>
            array(
                'id' => '9115',
                'province_id' => '91',
                'name' => 'KAB. WAROPEN',
            ),
            487 =>
            array(
                'id' => '9116',
                'province_id' => '91',
                'name' => 'KAB. BOVEN DIGOEL',
            ),
            488 =>
            array(
                'id' => '9117',
                'province_id' => '91',
                'name' => 'KAB. MAPPI',
            ),
            489 =>
            array(
                'id' => '9118',
                'province_id' => '91',
                'name' => 'KAB. ASMAT',
            ),
            490 =>
            array(
                'id' => '9119',
                'province_id' => '91',
                'name' => 'KAB. SUPIORI',
            ),
            491 =>
            array(
                'id' => '9120',
                'province_id' => '91',
                'name' => 'KAB. MAMBERAMO RAYA',
            ),
            492 =>
            array(
                'id' => '9121',
                'province_id' => '91',
                'name' => 'KAB. MAMBERAMO TENGAH',
            ),
            493 =>
            array(
                'id' => '9122',
                'province_id' => '91',
                'name' => 'KAB. YALIMO',
            ),
            494 =>
            array(
                'id' => '9123',
                'province_id' => '91',
                'name' => 'KAB. LANNY JAYA',
            ),
            495 =>
            array(
                'id' => '9124',
                'province_id' => '91',
                'name' => 'KAB. NDUGA',
            ),
            496 =>
            array(
                'id' => '9125',
                'province_id' => '91',
                'name' => 'KAB. PUNCAK',
            ),
            497 =>
            array(
                'id' => '9126',
                'province_id' => '91',
                'name' => 'KAB. DOGIYAI',
            ),
            498 =>
            array(
                'id' => '9127',
                'province_id' => '91',
                'name' => 'KAB. INTAN JAYA',
            ),
            499 =>
            array(
                'id' => '9128',
                'province_id' => '91',
                'name' => 'KAB. DEIYAI',
            ),
        ));
        \DB::table('regencies')->insert(array(
            0 =>
            array(
                'id' => '9171',
                'province_id' => '91',
                'name' => 'KOTA JAYAPURA',
            ),
            1 =>
            array(
                'id' => '9201',
                'province_id' => '92',
                'name' => 'KAB. SORONG',
            ),
            2 =>
            array(
                'id' => '9202',
                'province_id' => '92',
                'name' => 'KAB. MANOKWARI',
            ),
            3 =>
            array(
                'id' => '9203',
                'province_id' => '92',
                'name' => 'KAB. FAK FAK',
            ),
            4 =>
            array(
                'id' => '9204',
                'province_id' => '92',
                'name' => 'KAB. SORONG SELATAN',
            ),
            5 =>
            array(
                'id' => '9205',
                'province_id' => '92',
                'name' => 'KAB. RAJA AMPAT',
            ),
            6 =>
            array(
                'id' => '9206',
                'province_id' => '92',
                'name' => 'KAB. TELUK BINTUNI',
            ),
            7 =>
            array(
                'id' => '9207',
                'province_id' => '92',
                'name' => 'KAB. TELUK WONDAMA',
            ),
            8 =>
            array(
                'id' => '9208',
                'province_id' => '92',
                'name' => 'KAB. KAIMANA',
            ),
            9 =>
            array(
                'id' => '9209',
                'province_id' => '92',
                'name' => 'KAB. TAMBRAUW',
            ),
            10 =>
            array(
                'id' => '9210',
                'province_id' => '92',
                'name' => 'KAB. MAYBRAT',
            ),
            11 =>
            array(
                'id' => '9211',
                'province_id' => '92',
                'name' => 'KAB. MANOKWARI SELATAN',
            ),
            12 =>
            array(
                'id' => '9212',
                'province_id' => '92',
                'name' => 'KAB. PEGUNUNGAN ARFAK',
            ),
            13 =>
            array(
                'id' => '9271',
                'province_id' => '92',
                'name' => 'KOTA SORONG',
            ),
        ));
    }
}
