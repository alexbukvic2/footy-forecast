-- +goose Up
UPDATE teams SET short_name = 'BEL' WHERE name = 'Belgium';
UPDATE teams SET short_name = 'FRA' WHERE name = 'France';
UPDATE teams SET short_name = 'CRO' WHERE name = 'Croatia';
UPDATE teams SET short_name = 'SWE' WHERE name = 'Sweden';
UPDATE teams SET short_name = 'BRA' WHERE name = 'Brazil';
UPDATE teams SET short_name = 'URU' WHERE name = 'Uruguay';
UPDATE teams SET short_name = 'COL' WHERE name = 'Colombia';
UPDATE teams SET short_name = 'ESP' WHERE name = 'Spain';
UPDATE teams SET short_name = 'ENG' WHERE name = 'England';
UPDATE teams SET short_name = 'PAN' WHERE name = 'Panama';
UPDATE teams SET short_name = 'JPN' WHERE name = 'Japan';
UPDATE teams SET short_name = 'SEN' WHERE name = 'Senegal';
UPDATE teams SET short_name = 'SUI' WHERE name = 'Switzerland';
UPDATE teams SET short_name = 'MEX' WHERE name = 'Mexico';
UPDATE teams SET short_name = 'KOR' WHERE name = 'South Korea';
UPDATE teams SET short_name = 'AUS' WHERE name = 'Australia';
UPDATE teams SET short_name = 'IRN' WHERE name = 'Iran';
UPDATE teams SET short_name = 'KSA' WHERE name = 'Saudi Arabia';
UPDATE teams SET short_name = 'GER' WHERE name = 'Germany';
UPDATE teams SET short_name = 'ARG' WHERE name = 'Argentina';
UPDATE teams SET short_name = 'POR' WHERE name = 'Portugal';
UPDATE teams SET short_name = 'TUN' WHERE name = 'Tunisia';
UPDATE teams SET short_name = 'MAR' WHERE name = 'Morocco';
UPDATE teams SET short_name = 'EGY' WHERE name = 'Egypt';
UPDATE teams SET short_name = 'CZE' WHERE name = 'Czech Republic';
UPDATE teams SET short_name = 'AUT' WHERE name = 'Austria';
UPDATE teams SET short_name = 'TUR' WHERE name = 'Türkiye';
UPDATE teams SET short_name = 'NOR' WHERE name = 'Norway';
UPDATE teams SET short_name = 'SCO' WHERE name = 'Scotland';
UPDATE teams SET short_name = 'NED' WHERE name = 'Netherlands';
UPDATE teams SET short_name = 'CIV' WHERE name = 'Ivory Coast';
UPDATE teams SET short_name = 'GHA' WHERE name = 'Ghana';
UPDATE teams SET short_name = 'COD' WHERE name = 'Congo DR';
UPDATE teams SET short_name = 'RSA' WHERE name = 'South Africa';
UPDATE teams SET short_name = 'ALG' WHERE name = 'Algeria';
UPDATE teams SET short_name = 'JOR' WHERE name = 'Jordan';
UPDATE teams SET short_name = 'IRQ' WHERE name = 'Iraq';
UPDATE teams SET short_name = 'UZB' WHERE name = 'Uzbekistan';
UPDATE teams SET short_name = 'QAT' WHERE name = 'Qatar';
UPDATE teams SET short_name = 'PAR' WHERE name = 'Paraguay';
UPDATE teams SET short_name = 'ECU' WHERE name = 'Ecuador';
UPDATE teams SET short_name = 'USA' WHERE name = 'USA';
UPDATE teams SET short_name = 'HAI' WHERE name = 'Haiti';
UPDATE teams SET short_name = 'NZL' WHERE name = 'New Zealand';
UPDATE teams SET short_name = 'CAN' WHERE name = 'Canada';
UPDATE teams SET short_name = 'CUW' WHERE name = 'Curaçao';
UPDATE teams SET short_name = 'CPV' WHERE name = 'Cape Verde';
UPDATE teams SET short_name = 'BIH' WHERE name = 'Bosnia';

-- +goose Down
UPDATE teams SET short_name = ''
WHERE name IN (
    'Belgium', 'France', 'Croatia', 'Sweden', 'Brazil', 'Uruguay', 'Colombia',
    'Spain', 'England', 'Panama', 'Japan', 'Senegal', 'Switzerland', 'Mexico',
    'South Korea', 'Australia', 'Iran', 'Saudi Arabia', 'Germany', 'Argentina',
    'Portugal', 'Tunisia', 'Morocco', 'Egypt', 'Czech Republic', 'Austria',
    'Türkiye', 'Norway', 'Scotland', 'Netherlands', 'Ivory Coast', 'Ghana',
    'Congo DR', 'South Africa', 'Algeria', 'Jordan', 'Iraq', 'Uzbekistan',
    'Qatar', 'Paraguay', 'Ecuador', 'USA', 'Haiti', 'New Zealand', 'Canada',
    'Curaçao', 'Cape Verde', 'Bosnia'
);
