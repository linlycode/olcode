import copy from 'clipboard-copy'
export function copyToClipboard(text: string) {
	copy(text)
}
