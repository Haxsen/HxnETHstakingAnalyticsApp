'use client'

import { useState, useMemo } from 'react'

interface DonateModalProps {
  isOpen: boolean
  onClose: () => void
}

const WALLETS = [
  { name: 'BTC', address: process.env.NEXT_PUBLIC_WALLET_BTC || '' },
  { name: 'EVM(ETH/BNB/BASE/L2s)', address: process.env.NEXT_PUBLIC_WALLET_ETH || '' },
  { name: 'SOL', address: process.env.NEXT_PUBLIC_WALLET_SOL || '' },
  { name: 'SUI', address: process.env.NEXT_PUBLIC_WALLET_SUI || '' },
  { name: 'APTOS', address: process.env.NEXT_PUBLIC_WALLET_APTOS || '' },
  { name: 'NEAR', address: process.env.NEXT_PUBLIC_WALLET_NEAR || '' },
].filter(wallet => wallet.address !== '')

export default function DonateModal({ isOpen, onClose }: DonateModalProps) {
  const [copiedAddress, setCopiedAddress] = useState<string | null>(null)

  const copyToClipboard = async (address: string) => {
    try {
      await navigator.clipboard.writeText(address)
      setCopiedAddress(address)
      setTimeout(() => setCopiedAddress(null), 2000)
    } catch (err) {
      console.error('Failed to copy address:', err)
    }
  }

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div
        className="relative w-full max-w-md rounded-lg border p-6 shadow-xl"
        style={{
          backgroundColor: 'rgb(var(--bg-secondary))',
          borderColor: 'rgb(var(--border))',
        }}
      >
        {/* Header */}
        <div className="mb-4 flex items-center justify-between">
          <h3
            className="text-lg font-semibold"
            style={{ color: 'rgb(var(--text-primary))' }}
          >
            Support the Project
          </h3>
          <button
            onClick={onClose}
            className="text-xl hover:opacity-70 transition-opacity"
            style={{ color: 'rgb(var(--text-secondary))' }}
          >
            ×
          </button>
        </div>

        {/* Description */}
        <p
          className="mb-4 text-sm"
          style={{ color: 'rgb(var(--text-secondary))' }}
        >
          Your support helps keep this project running. Choose your preferred cryptocurrency:
        </p>

        {/* Wallet List */}
        <div className="space-y-3">
          {WALLETS.map((wallet) => (
            <div
              key={wallet.name}
              className="flex items-center justify-between rounded border p-3"
              style={{
                backgroundColor: 'rgb(var(--bg-primary))',
                borderColor: 'rgb(var(--border))',
              }}
            >
              <div className="flex-1 min-w-0">
                <div
                  className="font-medium text-sm"
                  style={{ color: 'rgb(var(--text-primary))' }}
                >
                  {wallet.name}
                </div>
                <div
                  className="text-xs font-mono break-all mt-1"
                  style={{ color: 'rgb(var(--text-secondary))' }}
                >
                  {wallet.address}
                </div>
              </div>
              <button
                onClick={() => copyToClipboard(wallet.address)}
                className="ml-3 px-3 py-1 text-xs rounded transition-colors"
                style={{
                  backgroundColor: copiedAddress === wallet.address ? '#10b981' : 'rgb(var(--border))',
                  color: copiedAddress === wallet.address ? 'white' : 'rgb(var(--text-primary))',
                }}
              >
                {copiedAddress === wallet.address ? 'Copied!' : 'Copy'}
              </button>
            </div>
          ))}
        </div>

        {/* Footer */}
        <div className="mt-4 pt-4 border-t" style={{ borderColor: 'rgb(var(--border))' }}>
          <p
            className="text-xs text-center"
            style={{ color: 'rgb(var(--text-secondary))' }}
          >
            Thank you for your support! ❤️
          </p>
        </div>
      </div>
    </div>
  )
}
