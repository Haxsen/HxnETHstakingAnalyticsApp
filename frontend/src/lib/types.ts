// API Response Types

export interface Token {
  id: number
  symbol: string
  name: string
  contract_address: string
  decimals: number
  blockchain: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface TokensResponse {
  tokens: Token[]
  count: number
}

export interface PricePoint {
  timestamp: number
  price: number
}

export interface TokenHistoryResponse {
  token_symbol: string
  price_history: PricePoint[]
  count: number
}

export interface TokenValuationResponse {
  // TODO: Define when valuation endpoint is implemented
  message: string
}

export interface CacheRefreshResponse {
  message: string
}

// Chart Data Types
export interface ChartDataPoint {
  timestamp: number
  date: string
  [tokenSymbol: string]: number | string
}

export interface ChartSeries {
  name: string
  type: 'line'
  data: (number | null)[]
  smooth: boolean
  symbol: 'none'
  lineStyle: {
    width: number
    type?: 'solid' | 'dashed'
  }
  itemStyle?: {
    color: string
  }
}

// UI State Types
export interface LoadingState {
  tokens: boolean
  history: boolean
  [key: string]: boolean
}

export interface ErrorState {
  tokens?: string
  history?: string
  [key: string]: string | undefined
}
