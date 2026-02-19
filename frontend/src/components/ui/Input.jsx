import React from 'react';

const Input = ({ label, type = 'text', value, onChange, placeholder = '', required = false, error = '', className = '', options = [], ...props }) => {
  const baseClassName = `w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${
    error ? 'border-red-500' : 'border-gray-300'
  }`;

  return (
    <div className={`mb-4 ${className}`}>
      {label && (
        <label className="block text-sm font-medium text-gray-700 mb-1">
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}
      {type === 'select' ? (
        <select
          value={value}
          onChange={onChange}
          required={required}
          className={baseClassName}
          {...props}
        >
          {options.length > 0 && !options.some(opt => opt.value === '') && (
            <option value="">-- Pilih --</option>
          )}
          {options.map((opt, idx) => (
            <option key={idx} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      ) : (
        <input
          type={type}
          value={value}
          onChange={onChange}
          placeholder={placeholder}
          required={required}
          className={baseClassName}
          {...props}
        />
      )}
      {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
    </div>
  );
};

export default Input;
