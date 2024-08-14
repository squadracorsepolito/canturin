export function getHexNumber(num: number) {
	return `0x${num.toString(16).padStart(2, '0')}`;
}
