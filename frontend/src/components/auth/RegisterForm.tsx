import React from 'react';
import { useForm } from 'react-hook-form';
import { useRegister } from '../../hooks/useAuth';
import Input from '../ui/Input';
import Button from '../ui/Button';
import type { RegisterRequest } from '../../types/user';

interface RegisterFormProps {
  onSuccess?: () => void;
}

const RegisterForm: React.FC<RegisterFormProps> = ({ onSuccess }) => {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterRequest & { confirmPassword: string }>();

  const { mutate: registerUser, isPending, error } = useRegister();

  const password = watch('password');

  const onSubmit = (data: RegisterRequest & { confirmPassword: string }) => {
    const { confirmPassword, ...registerData } = data;
    registerUser(registerData, {
      onSuccess: () => {
        onSuccess?.();
      },
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <Input
        label="Full Name"
        type="text"
        fullWidth
        error={errors.name?.message}
        {...register('name', {
          required: 'Name is required',
          minLength: {
            value: 2,
            message: 'Name must be at least 2 characters',
          },
        })}
      />

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
          minLength: {
            value: 8,
            message: 'Password must be at least 8 characters',
          },
          pattern: {
            value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
            message: 'Password must contain uppercase, lowercase, number, and special character',
          },
        })}
      />

      <Input
        label="Confirm Password"
        type="password"
        fullWidth
        error={errors.confirmPassword?.message}
        {...register('confirmPassword', {
          required: 'Please confirm your password',
          validate: (value) => value === password || 'Passwords do not match',
        })}
      />

      {error && (
        <div className="text-rose text-sm p-3 bg-rose/10 rounded-lg border border-rose/20">
          {error.message || 'Registration failed. Please try again.'}
        </div>
      )}

      <Button
        type="submit"
        variant="primary"
        fullWidth
        isLoading={isPending}
        disabled={isPending}
      >
        {isPending ? 'Creating Account...' : 'Create Account'}
      </Button>
    </form>
  );
};

export default RegisterForm;
