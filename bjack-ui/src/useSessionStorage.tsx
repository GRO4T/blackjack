import { useState, useEffect } from "react";

function getStorageValue<T>(key: string, defaultValue: T): T {
  const saved = sessionStorage.getItem(key);
  if (saved) {
    try {
      return JSON.parse(saved) as T;
    } catch {
      return defaultValue;
    }
  }
  return defaultValue;
}

export const useSessionStorage = <T,>(
  key: string,
  defaultValue: T,
): [T, React.Dispatch<React.SetStateAction<T>>] => {
  const [value, setValue] = useState<T>(() => {
    return getStorageValue(key, defaultValue);
  });

  useEffect(() => {
    sessionStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue];
};
