<script lang="ts">
	import type { SignalEnum } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import Divider from '$lib/components/divider/divider.svelte';
	import Readonly from '$lib/components/readonly/readonly.svelte';
	import { getSignalEnumState } from '$lib/panel/signal-enum/state.svelte';
	import type { PanelSectionProps } from '../types';

	let { entityId }: PanelSectionProps = $props();

	const ses = getSignalEnumState(entityId);
</script>

{#snippet section(signalEnum: SignalEnum)}
	<Attribute label="Size" desc="The current size of the enum in bits">
		<Readonly>
			<span class="font-medium">{signalEnum.size}</span>
		</Readonly>
	</Attribute>

	<Divider />

	<Attribute
		label="Minimum Size"
		desc="The minimum size of the enum in bits. It is useful when a signal should be of a fixed size"
	>
		<Readonly>
			<span class="font-medium">{signalEnum.minSize}</span>
		</Readonly>
	</Attribute>

	<Divider />

	<Attribute label="Maximum Index" desc="The maximum index of the values">
		<Readonly>
			<span class="font-medium">{signalEnum.maxIndex}</span>
		</Readonly>
	</Attribute>
{/snippet}

<secion>
	<h3 class="pb-5">Attributes</h3>

	{@render section(ses.entity)}
</secion>
