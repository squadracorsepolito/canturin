<!-- <script lang="ts">
	import { type SignalType } from '$lib/api/canturin';
	import { SignalTypeKind } from '$lib/api/github.com/squadracorsepolito/acmelib';
	import Summary from '$lib/components/summary/summary.svelte';
	import DescTextarea from '$lib/components/textarea/desc-textarea.svelte';
	import ReferenceTree from '$lib/components/tree/reference-tree.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import Panel from '../panel.svelte';
	import { createSignalTypeContext } from './signal-type-context.svelte';
	import SignalTypeAttributes from './signal-type-attributes.svelte';
	import SignalTypeHeading from './signal-type-heading.svelte';

	type Props = {
		entityId: string;
	};

	let { entityId }: Props = $props();

	let signalType = useSignalType(entityId);

	$effect(() => {
		signalType.reload(entityId);
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

		// if (sigType.size > 1) {
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
		// }

		// if (sigType.scale !== 1) {
		res.push({
			title: 'Scale',
			value: sigType.scale,
			desc: 'The scale value'
		});
		// }

		// if (sigType.offset !== 0) {
		res.push({
			title: 'Offset',
			value: sigType.offset,
			desc: 'The offset value'
		});
		// }

		return res;
	}

	function handleDesc(desc: string) {
		signalType.updateDesc(desc);
	}
</script>

{#snippet sigTypePanel(sigType: SignalType)}
	<SignalTypeHeading {entityId} />

	<section>
		<Summary infos={getSummaryInfos(sigType)} />
	</section>

	<SignalTypeAttributes {entityId} />

	{#if sigType.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(sigType.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#if !signalType.isLoading && signalType.entity}
		{@const _ = createSignalTypeContext(signalType.entity)}
		{@render sigTypePanel(signalType.entity)}
	{/if}
</Panel> -->
