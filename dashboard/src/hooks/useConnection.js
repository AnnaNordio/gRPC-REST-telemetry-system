import { useState, useEffect } from 'react';

export const useConnection = () => {
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    let socket;
    let reconnectTimeout;
    let isMounted = true;

    const connect = () => {
      const wsUrl = `ws://${window.location.hostname}:8080/ws`;
      socket = new WebSocket(wsUrl);
      socket.onopen = () => isMounted && setIsConnected(true);
      socket.onclose = () => {
        if (isMounted) {
          setIsConnected(false);
          reconnectTimeout = setTimeout(connect, 5000);
        }
      };
      socket.onerror = () => socket.close();
    };

    connect();
    return () => {
      isMounted = false;
      clearTimeout(reconnectTimeout);
      if (socket) socket.close();
    };
  }, []);

  return isConnected;
};