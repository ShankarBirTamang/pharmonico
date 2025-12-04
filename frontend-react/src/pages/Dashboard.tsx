import { Card, CardHeader, CardTitle, CardContent } from '../components/ui'

interface StatCardProps {
  title: string
  value: string | number
  change?: string
  trend?: 'up' | 'down' | 'neutral'
  icon: React.ReactNode
}

function StatCard({ title, value, change, trend, icon }: StatCardProps) {
  const trendColors = {
    up: 'text-emerald-600 bg-emerald-100',
    down: 'text-red-600 bg-red-100',
    neutral: 'text-slate-600 bg-slate-100',
  }

  return (
    <Card>
      <div className="flex items-start justify-between">
        <div>
          <p className="text-sm font-medium text-slate-500 dark:text-slate-400">{title}</p>
          <p className="mt-2 text-3xl font-bold text-slate-900 dark:text-white">{value}</p>
          {change && trend && (
            <p className={`mt-2 inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${trendColors[trend]}`}>
              {trend === 'up' && '↑ '}
              {trend === 'down' && '↓ '}
              {change}
            </p>
          )}
        </div>
        <div className="rounded-lg bg-primary-100 p-3 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400">
          {icon}
        </div>
      </div>
    </Card>
  )
}

export function Dashboard() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-slate-900 dark:text-white">Dashboard</h1>
        <p className="mt-1 text-slate-600 dark:text-slate-400">
          Overview of prescription routing activity
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Prescriptions"
          value="1,284"
          change="+12.5%"
          trend="up"
          icon={
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          }
        />
        <StatCard
          title="Pending Routing"
          value="42"
          change="-8.3%"
          trend="down"
          icon={
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
        />
        <StatCard
          title="Active Pharmacies"
          value="156"
          change="+3"
          trend="up"
          icon={
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          }
        />
        <StatCard
          title="Fulfillment Rate"
          value="94.2%"
          change="+2.1%"
          trend="up"
          icon={
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          }
        />
      </div>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Prescriptions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-slate-200 dark:border-slate-700">
                  <th className="pb-3 text-left text-sm font-semibold text-slate-900 dark:text-white">ID</th>
                  <th className="pb-3 text-left text-sm font-semibold text-slate-900 dark:text-white">Patient</th>
                  <th className="pb-3 text-left text-sm font-semibold text-slate-900 dark:text-white">Status</th>
                  <th className="pb-3 text-left text-sm font-semibold text-slate-900 dark:text-white">Pharmacy</th>
                  <th className="pb-3 text-left text-sm font-semibold text-slate-900 dark:text-white">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 dark:divide-slate-800">
                <tr>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">#RX-2024-001</td>
                  <td className="py-4 text-sm text-slate-900 dark:text-white">John Smith</td>
                  <td className="py-4">
                    <span className="inline-flex items-center rounded-full bg-emerald-100 px-2.5 py-0.5 text-xs font-medium text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400">
                      Fulfilled
                    </span>
                  </td>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">MedCare Pharmacy</td>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">Dec 4, 2024</td>
                </tr>
                <tr>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">#RX-2024-002</td>
                  <td className="py-4 text-sm text-slate-900 dark:text-white">Sarah Johnson</td>
                  <td className="py-4">
                    <span className="inline-flex items-center rounded-full bg-amber-100 px-2.5 py-0.5 text-xs font-medium text-amber-800 dark:bg-amber-900/30 dark:text-amber-400">
                      Awaiting Routing
                    </span>
                  </td>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">—</td>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">Dec 4, 2024</td>
                </tr>
                <tr>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">#RX-2024-003</td>
                  <td className="py-4 text-sm text-slate-900 dark:text-white">Michael Chen</td>
                  <td className="py-4">
                    <span className="inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-medium text-blue-800 dark:bg-blue-900/30 dark:text-blue-400">
                      Routed
                    </span>
                  </td>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">HealthFirst Rx</td>
                  <td className="py-4 text-sm text-slate-600 dark:text-slate-400">Dec 3, 2024</td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

