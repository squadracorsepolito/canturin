<script lang="ts">
	import type { SignalEnum } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import Readonly from '$lib/components/readonly/readonly.svelte';
	import { getSignalEnumState } from '$lib/state/signal-enum-state.svelte';
	import { text } from '../signal-enum/signal-enum-text';
	import type { PanelSectionProps } from '../types';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);
</script>

{#snippet section(signalEnum: SignalEnum)}
	<Attribute {...text.size}>
		<Readonly>
			<span class="font-medium">{signalEnum.size}</span>
		</Readonly>
	</Attribute>

	<Divider />

	<Attribute {...text.minSize}>
		<Readonly>
			<span class="font-medium">{signalEnum.minSize}</span>
		</Readonly>
	</Attribute>

	<Divider />

	<Attribute {...text.maxIndex}>
		<Readonly>
			<span class="font-medium">{signalEnum.maxIndex}</span>
		</Readonly>
	</Attribute>
{/snippet}

<secion>
	<h3 class="pb-5">{text.headings.attributes}</h3>

	{@render section(ses.entity)}
</secion>
