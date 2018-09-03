import log from 'src/infra/log'

export interface Codec {
	decode(rawStr: string): string
	encode(msg: string): string
}

class CodecImpl implements Codec {
	private header: string = "olcode_msg_header"

	public decode(rawStr: string): string {
		if (!rawStr.startsWith(this.header)) {
			log.error("invalid rawMsg: ", rawStr)
			return ""
		}

		return rawStr.substr(this.header.length)
	}

	public encode(msg: string): string {
		return `${this.header}${msg}`
	}
}

export default new CodecImpl()
