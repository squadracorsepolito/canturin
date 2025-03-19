<script lang="ts">
	import { SignalKind } from '$lib/api/canturin';

	type Props = {
		signalKind: SignalKind;
	};

	let { signalKind }: Props = $props();

	function getSignalKindString(kind: SignalKind) {
		switch (kind) {
			case SignalKind.SignalKindStandard:
				return 'Standard';
			case SignalKind.SignalKindEnum:
				return 'Enum';
			case SignalKind.SignalKindMultiplexer:
				return 'Multiplexer';
			default:
				return '';
		}
	}
</script>

<span
	class={[
		'badge badge-sm',
		signalKind === SignalKind.SignalKindStandard && 'badge-primary',
		signalKind === SignalKind.SignalKindEnum && 'badge-secondary',
		signalKind === SignalKind.SignalKindMultiplexer && 'badge-accent'
	]}>{getSignalKindString(signalKind)}</span
>
