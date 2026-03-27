export const formatTimestamp = (tsRaw) => {
  if (!tsRaw || tsRaw <= 0) return null;
  const ts = typeof tsRaw === 'object' ? tsRaw.toNumber?.() : Number(tsRaw);
  const date = new Date(ts / 1000); 

  return date.toLocaleTimeString('en-GB', {
    timeZone: 'UTC', 
    hour12: false, 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit',
    fractionalSecondDigits: 3 
  }).replace(/,/g, '');
};