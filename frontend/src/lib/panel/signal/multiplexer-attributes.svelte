<script lang="ts">
	import type { MultiplexerSignal } from '$lib/api/canturin';

	type Props = {
		entityId: string;
		signal: MultiplexerSignal;
	};

	let { entityId, signal }: Props = $props();
</script>

<pre>
    {JSON.stringify(signal, null, 2)}
</pre>
