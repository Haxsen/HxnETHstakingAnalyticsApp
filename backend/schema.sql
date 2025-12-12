-- Database schema for ETH Staking Analytics Backend

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
('pufETH', 'Puffer Staked ETH', '0xd9a442856c234a39a81a089c06451ebaa4306a72', 18),
('LSETH', 'Liquid Collective Staked ETH', '0x8c1bed5b9a0928467c9b1341da1d7bd5e10b6549', 18),
('RSETH', 'Kelp DAO Restaked ETH', '0xA1290d69c65A6Fe4DF752f95823fae25cB99e5A7', 18),
('METH', 'Mantle Staked Ether', '0xd5f7838f5c461feff7fe49ea5ebaf7728bb0adfa', 18),
('CBETH', 'Coinbase Wrapped Staked ETH', '0xBe9895146f7AF43049ca1c1AE358B0541Ea49704', 18),
('TETH', 'Treehouse ETH', '0xD11c452fc99cF405034ee446803b6F6c1F6d5ED8', 18),
('SFRXETH', 'Staked Frax Ether', '0xac3E018457B222d93114458476f3E3416Abbe38F', 18),
('CDCETH', 'Crypto.com Staked ETH', '0xfe18aE03741a5b84e39C295Ac9C856eD7991C38e', 18),
('UNIETH', 'Universal ETH', '0xF1376bceF0f78459C0Ed0ba5ddce976F1ddF51F4', 18)
ON CONFLICT (symbol) DO NOTHING;
