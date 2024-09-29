<script lang="ts">
	import { focusOnDisplay } from '$lib/actions';
	import { type SignalType } from '$lib/api/canturin';
	import { SignalTypeKind } from '$lib/api/github.com/squadracorsepolito/acmelib';
	import EditableNew from '$lib/components/editable/editable-new.svelte';
	import Editable from '$lib/components/editable/editable.svelte';
	import SignalTypeIcon from '$lib/components/icon/signal-type-icon.svelte';
	import AttributeInput from '$lib/components/input/attribute-input.svelte';
	import NameInput from '$lib/components/input/name-input.svelte';
	import NumberInput from '$lib/components/input/number-input.svelte';
	import Summary from '$lib/components/summary/summary.svelte';
	import DescTextarea from '$lib/components/textarea/desc-textarea.svelte';
	import ReferenceTree from '$lib/components/tree/reference-tree.svelte';
	import { useSignalType } from '$lib/state/signal-type-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import { z, ZodError } from 'zod';
	import Panel from './panel.svelte';

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

	function handleDesc(desc: string) {
		signalType.updateDesc(desc);
	}

	const val = z.number().min(1).max(64);
	let tmpSize = $state(signalType.entity ? signalType.entity.size : 16);
	let sizeErr = $derived.by(() => {
		try {
			val.parse(tmpSize);
		} catch (error) {
			if (error instanceof ZodError) {
				return error.issues[0].message;
			}
		}
	});

	const schema = z.object({
		size: z.number().min(1).max(64)
	});
</script>

{#snippet sigTypePanel(sigType: SignalType)}
	<section>
		<NameInput
			label="Signal Type Name"
			prefixName="signal_type"
			initialValue={sigType.name}
			onSubmit={(n) => {
				signalType.updateName(n);
			}}
			invalidNames={signalType.getInvalidNames()}
		>
			<SignalTypeIcon height="48" width="48" />
		</NameInput>

		<div class="pt-5">
			<DescTextarea initialDesc={sigType.desc} onSubmit={handleDesc} />
		</div>
	</section>

	<section>
		<Summary infos={getSummaryInfos(sigType)} />

		<!-- <div>
			<AttributeInput initialAttribute={sigType.size} onSubmit={(a) => console.log('got ', a)} />
		</div> -->
	</section>

	<section>
		<EditableNew
			{schema}
			initialValues={{ size: sigType.size }}
			onsubmit={() => console.log('submit')}
		>
			{#snippet placeholder()}
				<div class="stat place-items-center">
					<div class="stat-title">Size</div>

					<div class="stat-value">{sigType.size}</div>
				</div>
			{/snippet}

			{#snippet input({ fsm, values, errors })}
				<NumberInput
					bind:value={values.size}
					name="signal-type-size"
					errors={errors.size}
					focusOnDiplay
					onblur={() => fsm.send('BLUR')}
					onescape={() => fsm.send('ESCAPE')}
				/>
			{/snippet}
		</EditableNew>
	</section>

	{#if sigType.references}
		<section>
			<h4>References</h4>

			<ReferenceTree siblingNodes={getSignalReferenceTree(sigType.references)} depth={4} />
		</section>
	{/if}
{/snippet}

<Panel>
	{#if !signalType.isLoading && signalType.entity}
		{@render sigTypePanel(signalType.entity)}
	{/if}
</Panel>
