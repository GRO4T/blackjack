// https://blog.logrocket.com/using-localstorage-react-hooks/
import { useState, useEffect } from "react";

function getStorageValue(key: string, defaultValue: unknown) {
  const saved = sessionStorage.getItem(key);
  const initial = saved ? JSON.parse(saved) : defaultValue;
  return initial;
}

export const useSessionStorage = (key: string, defaultValue: unknown) => {
  const [value, setValue] = useState(() => {
    return getStorageValue(key, defaultValue);
  });

  useEffect(() => {
    sessionStorage.setItem(key, JSON.stringify(value));
  }, [key, value]);

  return [value, setValue];
};
