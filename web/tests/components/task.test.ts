import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import Task from '../../src/lib/components/Task.svelte';

describe('Task component', () => {
  it('should show the task name given name = "Test task"', () => {
    render(Task, { name: 'Test task' });

    const taskName = screen.getByText('Test task');

    expect(taskName.innerHTML).toBe('Test task');
  });

  it('should show the task name given name = "Another task"', () => {
    render(Task, { name: 'Another task' });

    const taskName = screen.getByText('Another task');

    expect(taskName.innerHTML).toBe('Another task');
  });
});