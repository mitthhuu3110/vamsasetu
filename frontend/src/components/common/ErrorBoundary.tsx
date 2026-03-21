import { Component, type ErrorInfo, type ReactNode } from 'react';
import Button from '../ui/Button';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo);
  }

  private handleReload = () => {
    window.location.reload();
  };

  public render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-ivory p-4">
          <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8 text-center">
            <div className="text-6xl mb-4">⚠️</div>
            <h1 className="text-2xl font-bold text-charcoal mb-2">
              Oops! Something went wrong
            </h1>
            <p className="text-charcoal/70 mb-6">
              We're sorry for the inconvenience. An unexpected error occurred.
            </p>
            {this.state.error && (
              <div className="bg-rose/10 border border-rose/20 rounded-lg p-4 mb-6 text-left">
                <p className="text-sm text-rose font-mono break-all">
                  {this.state.error.message}
                </p>
              </div>
            )}
            <Button
              variant="primary"
              onClick={this.handleReload}
              fullWidth
            >
              Reload Page
            </Button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
