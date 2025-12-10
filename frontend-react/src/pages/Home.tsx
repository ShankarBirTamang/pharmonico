import { Link } from 'react-router-dom'

export function Home() {
  return (
    <>
      <div className="text-center">
        <h1 className="text-4xl font-bold tracking-tight text-slate-900 sm:text-5xl md:text-6xl dark:text-white">
          Hello {' '}
          <span className="bg-gradient-to-r from-primary-600 to-teal-500 bg-clip-text text-transparent">
            Pharmonico
          </span>
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg text-slate-600 dark:text-slate-400">
          Specialty Prescription Routing Platform — Intelligent pharmacy matching, 
          real-time tracking, and seamless healthcare coordination.
        </p>
        
        <div className="mt-10 flex items-center justify-center gap-4">
          <Link to="/ops/dashboard" className="btn-primary">
            Get Started
          </Link>
          <button className="btn-secondary">
            Learn More
          </button>
        </div>
      </div>

      {/* Feature Cards */}
      <div className="mt-20 grid gap-8 md:grid-cols-3">
        <div className="card group hover:shadow-lg hover:border-primary-200 transition-all duration-300">
          <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-primary-100 text-primary-600 group-hover:bg-primary-600 group-hover:text-white transition-colors">
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-white">
            Prescription Intake
          </h3>
          <p className="mt-2 text-slate-600 dark:text-slate-400">
            NCPDP-compliant prescription parsing with intelligent validation and error handling.
          </p>
        </div>

        <div className="card group hover:shadow-lg hover:border-teal-200 transition-all duration-300">
          <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-teal-100 text-teal-600 group-hover:bg-teal-600 group-hover:text-white transition-colors">
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-white">
            Smart Routing
          </h3>
          <p className="mt-2 text-slate-600 dark:text-slate-400">
            AI-powered pharmacy matching based on specialty, location, and capacity scoring.
          </p>
        </div>

        <div className="card group hover:shadow-lg hover:border-amber-200 transition-all duration-300">
          <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 text-amber-600 group-hover:bg-amber-600 group-hover:text-white transition-colors">
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <h3 className="text-lg font-semibold text-slate-900 dark:text-white">
            Real-time Tracking
          </h3>
          <p className="mt-2 text-slate-600 dark:text-slate-400">
            Live status updates with event-driven architecture and comprehensive audit logging.
          </p>
        </div>
      </div>

      {/* Status Badge */}
      <div className="mt-16 flex justify-center">
        <div className="inline-flex items-center gap-2 rounded-full bg-emerald-100 px-4 py-2 text-sm font-medium text-emerald-800 dark:bg-emerald-900/30 dark:text-emerald-400">
          <span className="relative flex h-2 w-2">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75"></span>
            <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500"></span>
          </span>
          System Operational
        </div>
      </div>
    </>
  )
}

