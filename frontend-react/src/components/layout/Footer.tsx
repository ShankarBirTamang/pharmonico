export function Footer() {
  return (
    <footer className="border-t border-slate-200 bg-white dark:border-slate-800 dark:bg-slate-950">
      <div className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="flex flex-col items-center justify-between gap-4 sm:flex-row">
          <div className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-primary-500 to-teal-500">
              <span className="text-sm font-bold text-white">P</span>
            </div>
            <span className="font-semibold text-slate-900 dark:text-white">Pharmonico</span>
          </div>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            © {new Date().getFullYear()} Pharmonico. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  )
}

