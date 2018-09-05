import log from 'src/infra/log'

export interface Codec {
	decode(rawStr: string): string
	encode(msg: string): string
}

function b64EncodeUnicode(str: string) {
	// first we use encodeURIComponent to get percent-encoded UTF-8,
	// then we convert the percent encodings into raw bytes which
	// can be fed into btoa.
	return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g,
		function toSolidBytes(match, p1) {
			return String.fromCharCode(Number('0x' + p1))
		}))
}

function b64DecodeUnicode(str: string) {
	// Going backwards: from bytestream, to percent-encoding, to original string.
	return decodeURIComponent(atob(str).split('')
		.map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
		.join(''))
}

class CodecImpl implements Codec {
	private header: string = "olcode_msg_header"

	public decode(rawStr: string): string {
		const originStr = b64DecodeUnicode(rawStr)
		if (!originStr.startsWith(this.header)) {
			log.error("invalid rawMsg: ", originStr)
			return ""
		}

		return originStr.substr(this.header.length)
	}

	public encode(msg: string): string {
		return b64EncodeUnicode(`${this.header}${msg}`)
	}
}

export default new CodecImpl()
