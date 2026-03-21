import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import RegisterForm from '../components/auth/RegisterForm';
import Card, { CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '../components/ui/Card';

const RegisterPage: React.FC = () => {
  const navigate = useNavigate();

  const handleRegisterSuccess = () => {
    navigate('/dashboard');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-ivory p-4 relative overflow-hidden">
      {/* Rangoli background pattern */}
      <div className="absolute inset-0 opacity-5 rangoli-pattern" />

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md relative z-10"
      >
        <div className="text-center mb-8">
          <motion.h1
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2, duration: 0.5 }}
            className="font-display text-4xl font-bold text-charcoal mb-2"
          >
            VamsaSetu
          </motion.h1>
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3, duration: 0.5 }}
            className="text-charcoal/70"
          >
            Connect your family, preserve your heritage
          </motion.p>
        </div>

        <Card className="shadow-xl border-saffron/20">
          <CardHeader>
            <CardTitle className="text-2xl text-center">Create Account</CardTitle>
            <CardDescription className="text-center">
              Join VamsaSetu to start building your family tree
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RegisterForm onSuccess={handleRegisterSuccess} />
          </CardContent>
          <CardFooter className="flex flex-col items-center space-y-2">
            <p className="text-sm text-charcoal/60">
              Already have an account?{' '}
              <Link
                to="/login"
                className="text-saffron hover:text-saffron/80 font-medium transition-colors"
              >
                Sign in
              </Link>
            </p>
          </CardFooter>
        </Card>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5, duration: 0.5 }}
          className="mt-6 text-center text-sm text-charcoal/50"
        >
          <p>© 2026 VamsaSetu. Connecting generations.</p>
        </motion.div>
      </motion.div>
    </div>
  );
};

export default RegisterPage;
