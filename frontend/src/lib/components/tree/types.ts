import type { Component } from 'svelte';
import type { IconProps } from '../icon/types';

export type TreeNode = {
	name: string;
	icon: Component<IconProps>;
	childNodes: TreeNode[];
	onclick?: () => void;
};

export type ReferenceTreeNode = {
	name: string;
	childNodes: ReferenceTreeNode[];
};
