import React from 'react';
import { Loader2, CheckCircle2, XCircle, Clock } from 'lucide-react';

export type ScrapingStatus = 'pending' | 'running' | 'completed' | 'failed';

interface ScrapingStatusProps {
    status: ScrapingStatus;
    message?: string;
    resultsCount?: number;
    error?: string;
    onRetry?: () => void;
}

export const ScrapingStatusDisplay: React.FC<ScrapingStatusProps> = ({
    status,
    message,
    resultsCount,
    error,
    onRetry,
}) => {
    const getStatusConfig = () => {
        switch (status) {
            case 'pending':
                return {
                    icon: <Clock className="w-6 h-6 text-blue-500" />,
                    title: 'Queued',
                    description: message || 'Waiting to start scraping...',
                    bgColor: 'bg-blue-50',
                    borderColor: 'border-blue-200',
                    textColor: 'text-blue-700',
                };
            case 'running':
                return {
                    icon: <Loader2 className="w-6 h-6 text-blue-600 animate-spin" />,
                    title: 'Scraping Tabelog',
                    description: message || 'Searching for restaurants...',
                    bgColor: 'bg-blue-50',
                    borderColor: 'border-blue-300',
                    textColor: 'text-blue-700',
                };
            case 'completed':
                return {
                    icon: <CheckCircle2 className="w-6 h-6 text-green-600" />,
                    title: 'Completed',
                    description: message || `Found ${resultsCount || 0} restaurant${resultsCount === 1 ? '' : 's'}`,
                    bgColor: 'bg-green-50',
                    borderColor: 'border-green-200',
                    textColor: 'text-green-700',
                };
            case 'failed':
                return {
                    icon: <XCircle className="w-6 h-6 text-red-600" />,
                    title: 'Failed',
                    description: error || message || 'An error occurred while scraping',
                    bgColor: 'bg-red-50',
                    borderColor: 'border-red-200',
                    textColor: 'text-red-700',
                };
        }
    };

    const config = getStatusConfig();

    return (
        <div
            className={`
        ${config.bgColor} ${config.borderColor} ${config.textColor}
        border-2 rounded-lg p-4 mb-4
        transition-all duration-300 ease-in-out
        animate-fadeIn
      `}
        >
            <div className="flex items-start gap-3">
                <div className="flex-shrink-0 mt-0.5">
                    {config.icon}
                </div>
                <div className="flex-1 min-w-0">
                    <h4 className="font-semibold text-sm mb-1">
                        {config.title}
                    </h4>
                    <p className="text-sm opacity-90">
                        {config.description}
                    </p>
                    {status === 'failed' && onRetry && (
                        <button
                            onClick={onRetry}
                            className="mt-3 px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 transition-colors"
                        >
                            Retry
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};
