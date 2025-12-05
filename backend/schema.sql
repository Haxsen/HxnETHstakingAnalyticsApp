-- Database schema for ETH Staking Analytics Backend
-- Run this against your PostgreSQL database

CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    contract_address VARCHAR(42) NOT NULL UNIQUE,
    decimals INTEGER NOT NULL DEFAULT 18,
    blockchain VARCHAR(20) NOT NULL DEFAULT 'ethereum',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial LST tokens
INSERT INTO tokens (symbol, name, contract_address, decimals) VALUES
('wstETH', 'Wrapped Lido Staked Ether', '0x7f39c581f595b53c5cb19bd0b3f8da6c935e2ca0', 18),
('ankrETH', 'Ankr Staked ETH', '0xe95a203b1a91a908f9b9ce46459d101078c2c3cb', 18),
('rETH', 'Rocket Pool ETH', '0xae78736cd615f374d3085123a210448e74fc6393', 18),
('wBETH', 'Wrapped Binance Beacon ETH', '0xa2e3356610840701bdf5611a53974510ae27e2e1', 18),
('pufETH', 'Puffer Staked ETH', '0xd9a442856c234a39a81a089c06451ebaa4306a72', 18)
ON CONFLICT (symbol) DO NOTHING;
