'use client'

import { useState, useEffect } from 'react'
import Chart from '@/components/Chart'
import ValuationTable from '@/components/ValuationTable'
import ThemeToggle from '@/components/ThemeToggle'
import DonateButton from '@/components/DonateButton'
import { Token, ValuationData, LoadingState, ErrorState } from '@/lib/types'
import { fetchTokens, fetchTokenHistory, fetchValuations, getErrorMessage } from '@/lib/api'

export default function Home() {
  const [tokens, setTokens] = useState<Token[]>([])
  const [histories, setHistories] = useState<Record<string, any>>({})
  const [valuations, setValuations] = useState<ValuationData[]>([])
  const [loading, setLoading] = useState<LoadingState>({
    tokens: true,
    history: false,
    valuations: false,
  })
  const [errors, setErrors] = useState<ErrorState>({})

  // Fetch tokens on mount
  useEffect(() => {
    const loadTokens = async () => {
      try {
        setLoading(prev => ({ ...prev, tokens: true }))
        setErrors(prev => ({ ...prev, tokens: undefined }))

        const response = await fetchTokens()
        setTokens(response.tokens)
      } catch (error) {
        setErrors(prev => ({ ...prev, tokens: getErrorMessage(error) }))
      } finally {
        setLoading(prev => ({ ...prev, tokens: false }))
      }
    }

    loadTokens()
  }, [])

  // Fetch price histories when tokens are loaded
  useEffect(() => {
    if (tokens.length === 0) return

    const loadHistories = async () => {
      setLoading(prev => ({ ...prev, history: true }))
      setErrors(prev => ({ ...prev, history: undefined }))

      const newHistories: Record<string, any> = {}

      try {
        // Fetch history for each token
        const promises = tokens.map(async (token) => {
          try {
            const response = await fetchTokenHistory(token.symbol)
            newHistories[token.symbol] = response
          } catch (error) {
            console.error(`Failed to fetch history for ${token.symbol}:`, error)
            // Continue with other tokens even if one fails
          }
        })

        await Promise.all(promises)
        setHistories(newHistories)
      } catch (error) {
        setErrors(prev => ({ ...prev, history: getErrorMessage(error) }))
      } finally {
        setLoading(prev => ({ ...prev, history: false }))
      }
    }

    loadHistories()
  }, [tokens])

  // Fetch valuations when tokens are loaded
  useEffect(() => {
    if (tokens.length === 0) return

    const loadValuations = async () => {
      try {
        setLoading(prev => ({ ...prev, valuations: true }))
        setErrors(prev => ({ ...prev, valuations: undefined }))

        const response = await fetchValuations()
        setValuations(response.valuations)
      } catch (error) {
        setErrors(prev => ({ ...prev, valuations: getErrorMessage(error) }))
      } finally {
        setLoading(prev => ({ ...prev, valuations: false }))
      }
    }

    loadValuations()
  }, [tokens])

  const isLoading = loading.tokens || loading.history
  const hasError = errors.tokens || errors.history || errors.valuations

  return (
    <main className="min-h-screen py-8 transition-colors duration-200" style={{ backgroundColor: 'rgb(var(--bg-primary))' }}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8 flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2" style={{ color: 'rgb(var(--text-primary))' }}>
              ETH Liquid Staking Analytics Dashboard
            </h1>
          </div>
          <div className="flex items-center gap-3">
            <DonateButton />
            <ThemeToggle />
          </div>
        </div>

        {/* LST Explanation */}
        <div className="mb-8 p-4 rounded-lg border" style={{
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderColor: 'rgb(59, 130, 246)',
        }}>
          <h3 className="text-lg font-semibold mb-2" style={{ color: 'rgb(var(--text-primary))' }}>
            What are Liquid Staking Tokens (LSTs)?
          </h3>
          <p style={{ color: 'rgb(var(--text-secondary))' }}>
            Liquid Staked Tokens represent your staked assets (like ETH) that continue to earn staking rewards while remaining fully liquid and tradeable.
            Unlike traditional staking which locks your assets, LSTs allow you to earn yield without sacrificing liquidity. However, you are trusting the underlying protocol
            and its security measures.
          </p>
        </div>

        {/* Error Display */}
        {hasError && (
          <div className="mb-6 p-4 rounded-md border" style={{
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            borderColor: 'rgb(220, 38, 38)',
          }}>
            <div className="flex">
              <div className="ml-3">
                <h3 className="text-sm font-medium" style={{ color: 'rgb(185, 28, 28)' }}>
                  Error loading data
                </h3>
                <div className="mt-2 text-sm" style={{ color: 'rgb(153, 27, 27)' }}>
                  {errors.tokens && <p>Tokens: {errors.tokens}</p>}
                  {errors.history && <p>Price data: {errors.history}</p>}
                  {errors.valuations && <p>Valuation data: {errors.valuations}</p>}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Chart Container */}
        <div className="rounded-lg shadow-sm border p-6 transition-colors duration-200" style={{
          backgroundColor: 'rgb(var(--bg-secondary))',
          borderColor: 'rgb(var(--border))'
        }}>
          {isLoading ? (
            <div className="flex items-center justify-center h-96">
              <div className="text-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 mx-auto mb-4" style={{ borderBottomColor: 'rgb(59, 130, 246)' }}></div>
                <p style={{ color: 'rgb(var(--text-secondary))' }}>Loading price data (first time may take ~3 minutes)...</p>
              </div>
            </div>
          ) : tokens.length > 0 ? (
            <Chart
              tokens={tokens}
              histories={histories}
              loading={false}
            />
          ) : (
            <div className="flex items-center justify-center h-96">
              <div className="text-center">
                <p style={{ color: 'rgb(var(--text-secondary))' }}>No token data available</p>
              </div>
            </div>
          )}
        </div>

        {/* Valuation Table */}
        <div className="mt-8">
          {/* DYOR Note */}
          <div className="mb-4 p-3 rounded-md border" style={{
            backgroundColor: 'rgba(245, 158, 11, 0.1)',
            borderColor: 'rgb(245, 158, 11)',
          }}>
            <p className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>
              <strong>⚠️ DYOR:</strong> These valuations are based on 1-year historical data and should not be considered financial advice.
              Always do your own research before making investment decisions.
            </p>
          </div>

          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold" style={{ color: 'rgb(var(--text-primary))' }}>
              Valuation Metrics
            </h2>
            {valuations.length > 0 && (
              <p className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>
                Last updated: {new Date(valuations[0]?.last_updated || '').toLocaleString()}
              </p>
            )}
          </div>
          <div className="rounded-lg shadow-sm border p-6 transition-colors duration-200" style={{
            backgroundColor: 'rgb(var(--bg-secondary))',
            borderColor: 'rgb(var(--border))'
          }}>
            <ValuationTable
              valuations={valuations}
              loading={loading.valuations}
            />
          </div>

          {/* Risk Warning */}
          <div className="mt-4 p-3 rounded-md border" style={{
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            borderColor: 'rgb(239, 68, 68)',
          }}>
            <p className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>
              <strong>⚠️ Smart Contract Risk:</strong> When investing in LSTs, you are trusting the underlying protocol and its security measures.
              Research the company thoroughly - any protocol can suffer hacks or exploits. For example, StakeWise (providing OSETH) suffered an attack in early November 2025.
              Always understand the risks before investing.
            </p>
          </div>
        </div>

        {/* Token List */}
        {tokens.length > 0 && (
          <div className="mt-8">
            <h2 className="text-xl font-semibold mb-4" style={{ color: 'rgb(var(--text-primary))' }}>
              Tracked Tokens ({tokens.length})
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {tokens.map((token) => (
                <div
                  key={token.symbol}
                  className="rounded-lg shadow-sm border p-4 transition-colors duration-200"
                  style={{
                    backgroundColor: 'rgb(var(--bg-secondary))',
                    borderColor: 'rgb(var(--border))'
                  }}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="font-medium" style={{ color: 'rgb(var(--text-primary))' }}>{token.symbol}</h3>
                      <p className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>{token.name}</p>
                    </div>
                    <div className="text-right">
                      <span className="text-xs uppercase tracking-wide" style={{ color: 'rgb(var(--text-secondary))' }}>
                        {token.blockchain}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Author Note */}
        <div className="mt-12 text-center">
          <div className="p-4 rounded-lg border" style={{
            backgroundColor: 'rgba(59, 130, 246, 0.05)',
            borderColor: 'rgb(59, 130, 246)',
          }}>
            <p className="text-sm" style={{ color: 'rgb(var(--text-secondary))' }}>
              Made by <strong>Haxsen</strong> - A simple dashboard to share knowledge about Liquid Staking Tokens.
              <br />
              Contact: haxsenmail@gmail.com | Website: <a href="https://haxsen.github.io" target="_blank" rel="noopener noreferrer" className="underline hover:no-underline" style={{ color: 'rgb(59, 130, 246)' }}>haxsen.github.io</a>
            </p>
          </div>
        </div>
      </div>
    </main>
  )
}
