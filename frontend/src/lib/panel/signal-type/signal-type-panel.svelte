<script lang="ts">
	import { type SignalType } from '$lib/api/canturin';
	import { SignalTypeKind } from '$lib/api/github.com/squadracorsepolito/acmelib';
	import EditableNew from '$lib/components/editable/editable-new.svelte';
	import SignalTypeIcon from '$lib/components/icon/signal-type-icon.svelte';
	import NameInput from '$lib/components/input/name-input.svelte';
	import NumberInput from '$lib/components/input/number-input.svelte';
	import Summary from '$lib/components/summary/summary.svelte';
	import DescTextarea from '$lib/components/textarea/desc-textarea.svelte';
	import ReferenceTree from '$lib/components/tree/reference-tree.svelte';
	import { useSignalType } from '$lib/state/signal-type-state.svelte';
	import { getSignalReferenceTree } from '$lib/utils';
	import { z, ZodError } from 'zod';
	import Panel from '../panel.svelte';
	import { AttributePlaceholder } from '$lib/components/placeholder';
	import { maxSchema, minSchema, nameSchema, sizeSchema } from './schema';
	import { UpdateName } from '$lib/api/canturin/signaltypeservice';
	import TextInput from '$lib/components/input/text-input.svelte';

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
</script>

{#snippet attributes(sigType: SignalType)}
	<div class="flex-1">
		<EditableNew
			schema={sizeSchema}
			initialValues={{ size: sigType.size }}
			onsubmit={() => console.log('submit')}
		>
			{#snippet placeholder()}
				<AttributePlaceholder label="Size" value={sigType.size} />
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
	</div>

	<div class="flex-1">
		<EditableNew
			schema={minSchema}
			initialValues={{ min: sigType.min }}
			onsubmit={() => console.log('submit')}
		>
			{#snippet placeholder()}
				<AttributePlaceholder label="Min" value={sigType.min} />
			{/snippet}

			{#snippet input({ fsm, values, errors })}
				<NumberInput
					bind:value={values.min}
					name="signal-type-min"
					errors={errors.min}
					focusOnDiplay
					onblur={() => fsm.send('BLUR')}
					onescape={() => fsm.send('ESCAPE')}
					label="Min"
				/>
			{/snippet}
		</EditableNew>
	</div>

	<div class="flex-1">
		<EditableNew
			schema={maxSchema}
			initialValues={{ max: sigType.max }}
			onsubmit={() => console.log('submit')}
		>
			{#snippet placeholder()}
				<AttributePlaceholder label="Max" value={sigType.max} />
			{/snippet}

			{#snippet input({ fsm, values, errors })}
				<NumberInput
					bind:value={values.max}
					name="signal-type-max"
					errors={errors.max}
					focusOnDiplay
					onblur={() => fsm.send('BLUR')}
					onescape={() => fsm.send('ESCAPE')}
					label="Max"
				/>
			{/snippet}
		</EditableNew>
	</div>
{/snippet}

{#snippet sigTypePanel(sigType: SignalType)}
	<section>
		<!-- <NameInput
			label="Signal Type Name"
			prefixName="signal_type"
			initialValue={sigType.name}
			onSubmit={(n) => {
				signalType.updateName(n);
			}}
			invalidNames={signalType.getInvalidNames()}
		>
			<SignalTypeIcon height="48" width="48" />
		</NameInput> -->

		<div class="flex">
			<SignalTypeIcon height="48" width="48" />

			<EditableNew
				schema={nameSchema(signalType.getInvalidNames())}
				hidePlaceholder
				initialValues={{ name: sigType.name }}
				onsubmit={({ name }) => {
					signalType.updateName(name);
				}}
			>
				{#snippet placeholder()}
					<div>{sigType.name}</div>
				{/snippet}

				{#snippet input({ fsm, values, errors })}
					<TextInput
						label="Name"
						name="signal-type-name"
						bind:value={values.name}
						errors={errors.name}
						focusOnDiplay
						onblur={() => fsm.send('BLUR')}
						onescape={() => fsm.send('ESCAPE')}
					/>
				{/snippet}
			</EditableNew>
		</div>

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

	<section class="flex gap-3">
		{@render attributes(sigType)}
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
