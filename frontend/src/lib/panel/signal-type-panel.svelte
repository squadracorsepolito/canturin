<script lang="ts">
	import { type SignalType } from '$lib/api/canturin';
	import { SignalTypeKind } from '$lib/api/github.com/squadracorsepolito/acmelib';
	import Summary from '$lib/components/summary/summary.svelte';
	import ReferenceTree from '$lib/components/tree/reference-tree.svelte';
	import { useSignalType } from '$lib/state/signal-type-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import Panel from './panel.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let state = useSignalType(entityId);

	$effect(() => {
		state.reload(entityId);
	});

	function getKindString(kind: SignalTypeKind) {
		switch (kind) {
			case SignalTypeKind.SignalTypeKindFlag:
				return 'Flag';
			case SignalTypeKind.SignalTypeKindInteger:
				return 'Integer';
			case SignalTypeKind.SignalTypeKindDecimal:
				return 'Decimal';
			case SignalTypeKind.SignalTypeKindCustom:
				return 'Custom';
			default:
				return '';
		}
	}

	function getSummaryInfos(sigType: SignalType) {
		const res = [
			{
				title: 'Kind',
				value: getKindString(sigType.kind),
				desc: 'The of the type'
			},
			{
				title: 'Size',
				value: sigType.size,
				desc: 'The size in bits'
			}
		];

		if (sigType.size > 1) {
			res.push({
				title: 'Min',
				value: sigType.min,
				desc: 'The minimum value'
			});
			res.push({
				title: 'Max',
				value: sigType.max,
				desc: 'The maximum value'
			});
		}

		if (sigType.scale !== 1) {
			res.push({
				title: 'Scale',
				value: sigType.scale,
				desc: 'The scale value'
			});
		}

		if (sigType.offset !== 0) {
			res.push({
				title: 'Offset',
				value: sigType.offset,
				desc: 'The offset value'
			});
		}

		return res;
	}
</script>

{#snippet sigTypePanel(sigType: SignalType)}
	<section>
		<h3>{sigType.name}</h3>
		<p>{sigType.desc}</p>
	</section>

	<section>
		<Summary infos={getSummaryInfos(sigType)} />
	</section>

	{#if sigType.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(sigType.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#if !state.isLoading && state.signalType}
		{@render sigTypePanel(state.signalType)}
	{/if}
</Panel>
