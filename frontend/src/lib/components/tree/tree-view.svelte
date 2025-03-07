<script lang="ts" generics="T extends {[K in keyof T]: any } & {children?: T[] | null} ">
	import { uniqueId, type KeyOfString } from '$lib/utils';
	import { mergeProps, normalizeProps, useMachine } from '@zag-js/svelte';
	import * as tree from '@zag-js/tree-view';
	import { AltArrowIcon } from '../icon';
	import { type Component, type Snippet } from 'svelte';
	import type { IconProps } from '../icon/types';

	type Props = {
		root: T;
		selectedValue: string;
		valueKey: KeyOfString<T>;
		labelKey: KeyOfString<T>;
		getIcon: (node: T) => Component<IconProps>;
		onrootclick: () => void;
		onselect?: (value: string) => void;
		ondelete?: (value: string) => void;
		actions: Snippet<[{ collapse: () => void }]>;
	};

	let {
		root,
		selectedValue = $bindable(),
		valueKey,
		labelKey,
		getIcon,
		onrootclick,
		onselect,
		ondelete,
		actions
	}: Props = $props();

	const collection = tree.collection({
		rootNode: root,
		nodeToValue: (node) => node[valueKey],
		nodeToString: (node) => node[labelKey]
	});

	const zagTreeProps: tree.Props = $derived({
		id: uniqueId(),
		collection,
		selectedValue: [selectedValue],
		onSelectionChange: (details) => {
			selectedValue = details.selectedValue[0];
			onselect?.(details.selectedValue[0]);
		}
	});

	const service = useMachine(tree.machine, () => zagTreeProps);

	const api = $derived(tree.connect(service, normalizeProps));

	const RootIcon = getIcon(root);

	const treeProps = $derived(
		mergeProps(api.getTreeProps(), {
			onkeydown: (event: KeyboardEvent) => {
				if (!selectedValue) {
					return;
				}
				if (event.key === 'Delete' || (event.metaKey && event.key === 'Backspace')) {
					ondelete?.(selectedValue);
					selectedValue = '';
				}
			}
		})
	);

	const treeRootLableProps = $derived(
		mergeProps(api.getLabelProps(), {
			onclick: () => {
				onrootclick();
			}
		})
	);
</script>

{#snippet treeNode(node: T, indexPath: number[])}
	{@const nodeState = api.getNodeState({ indexPath, node })}
	{@const Icon = getIcon(node)}

	{#if nodeState.isBranch}
		<div {...api.getBranchProps({ node, indexPath })}>
			<div {...api.getBranchControlProps({ node, indexPath })}>
				<span {...api.getBranchIndicatorProps({ node, indexPath })}>
					<AltArrowIcon height="16" width="16" />
				</span>

				<span>
					<Icon height="16" width="16" />
				</span>

				<span {...api.getBranchTextProps({ node, indexPath })}>
					{node[labelKey]}
				</span>
			</div>

			<div {...api.getBranchContentProps({ node, indexPath })}>
				<div {...api.getBranchIndentGuideProps({ node, indexPath })}>
					<span></span>
				</div>

				{#if node.children}
					{#each node.children as child, idx}
						{@render treeNode(child, [...indexPath, idx])}
					{/each}
				{/if}
			</div>
		</div>
	{:else}
		<div {...api.getItemProps({ node, indexPath })}>
			<span>
				<Icon height="16" width="16" />
			</span>

			<span>
				{node[labelKey]}
			</span>
		</div>
	{/if}
{/snippet}

<div {...api.getRootProps()} class="overflow-x-hidden">
	<div class="flex items-center gap-2 p-2 pt-3">
		<span>
			<RootIcon height="20" width="20" />
		</span>

		<button {...treeRootLableProps}>
			{root[labelKey]}
		</button>

		{@render actions({ collapse: api.collapse })}
	</div>

	<div {...treeProps}>
		{#if collection.rootNode.children}
			{#each collection.rootNode.children as child, idx}
				{@render treeNode(child, [idx])}
			{/each}
		{/if}
	</div>
</div>

<style lang="postcss">
	[data-scope='tree-view'] [data-part='label'] {
		@apply font-medium text-sm truncate flex-1 text-left;
	}

	[data-scope='tree-view'] [data-part='branch-control'],
	[data-scope='tree-view'] [data-part='item'] {
		@apply relative flex items-center py-1 gap-2 border border-transparent outline-none;

		&:hover {
			@apply bg-base-content/20;
		}

		&:focus,
		&[data-selected] {
			@apply bg-primary-ghost text-primary;
		}

		&:focus {
			@apply border-primary;
		}
	}

	[data-scope='tree-view'] [data-part='branch-control'] {
		padding-inline-start: calc(var(--depth) * 0.5rem);
		padding-inline-end: 0.5rem;
	}

	[data-scope='tree-view'] [data-part='item'] {
		@apply cursor-pointer;

		padding-inline-start: calc(var(--depth) * 0.5rem + 1rem);
		padding-inline-end: 0.5rem;

		span:nth-child(2) {
			@apply text-xs truncate;
		}
	}

	[data-scope='tree-view'] [data-part='branch-indicator'] {
		&[data-state='closed'] {
			@apply -rotate-90;
		}
	}

	[data-scope='tree-view'] [data-part='branch-content'] {
		@apply relative;
	}

	[data-scope='tree-view'] [data-part='branch-text'] {
		@apply text-xs truncate;
	}

	[data-scope='tree-view'] [data-part='branch-indent-guide'] {
		@apply absolute h-full w-[0.125rem] py-1;

		left: calc(var(--depth) * 0.5rem + 0.375rem);

		span {
			@apply block h-full w-full bg-base-content/20 rounded-btn;
		}
	}
</style>
