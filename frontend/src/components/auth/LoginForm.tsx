import React from 'react';
import { useForm } from 'react-hook-form';
import { useLogin } from '../../hooks/useAuth';
import Input from '../ui/Input';
import Button from '../ui/Button';
import type { LoginRequest } from '../../types/user';

interface LoginFormProps {
  onSuccess?: () => void;
}

const LoginForm: React.FC<LoginFormProps> = ({ onSuccess }) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginRequest>();

  const { mutate: login, isPending, error } = useLogin();

  const onSubmit = (data: LoginRequest) => {
    login(data, {
      onSuccess: () => {
        onSuccess?.();
      },
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <Input
        label="Email"
        type="email"
        fullWidth
        error={errors.email?.message}
        {...register('email', {
          required: 'Email is required',
          pattern: {
            value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
            message: 'Invalid email address',
          },
        })}
      />

      <Input
        label="Password"
        type="password"
        fullWidth
        error={errors.password?.message}
        {...register('password', {
          required: 'Password is required',
        })}
      />

      {error && (
        <div className="text-rose text-sm p-3 bg-rose/10 rounded-lg border border-rose/20">
          {error.message || 'Login failed. Please try again.'}
        </div>
      )}

      <Button
        type="submit"
        variant="primary"
        fullWidth
        isLoading={isPending}
        disabled={isPending}
      >
        {isPending ? 'Logging in...' : 'Login'}
      </Button>
    </form>
  );
};

export default LoginForm;
