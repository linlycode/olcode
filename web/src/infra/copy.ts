export function copyToClipboard(text: string) {
	const input = document.createElement('input')
	// TODO: this input may cause dom change
	input.setAttribute('value', text)
	document.body.appendChild(input)
	input.select()
	document.execCommand('copy')
	document.body.removeChild(input)
}
