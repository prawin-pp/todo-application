import { fireEvent, render, screen } from '@testing-library/svelte';
import { describe, expect, it, vi } from 'vitest';
import Task from '../../src/lib/components/Task.svelte';

describe('Task component', () => {
  it('should show the task name given name = "Test task"', () => {
    const task: Task = { name: 'Test task', description: '', dueDate: '', completed: false };
    render(Task, { task });

    const taskName = screen.getByText('Test task');

    expect(taskName?.innerHTML).toBe('Test task');
  });

  it('should show the task name given name = "Another task"', () => {
    const task: Task = { name: 'Another task', description: '', dueDate: '', completed: false };
    render(Task, { task });

    const taskName = screen.getByText('Another task');

    expect(taskName?.innerHTML).toBe('Another task');
  });

  it('should not show the task name given not set name', () => {
    const task: Task = { name: '', description: '', dueDate: '', completed: false };
    render(Task, { task });

    const taskName = screen.queryByText('undefined');

    expect(taskName).toBeNull();
  });

  it('should show the task description given description = "Test description"', () => {
    const task: Task = { name: '', description: 'Test description', dueDate: '', completed: false };
    render(Task, { task });

    const taskDescription = screen.getByText('Test description');

    expect(taskDescription.innerHTML).toBe('Test description');
  });

  it('should not show the task description when not set description', () => {
    const task: Task = { name: '', description: '', dueDate: '', completed: false };
    render(Task, { task });

    const taskDescription = screen.queryByText('undefined');

    expect(taskDescription).toBeNull();
  });

  it('should show the task due date in format "DD MMM YYYY" given dueDate = "2021-10-10"', () => {
    const task: Task = { name: '', description: '', dueDate: '2021-10-10', completed: false };
    render(Task, { task });

    const taskDueDate = screen.queryByText('10 Oct 2021');

    expect(taskDueDate).not.toBeNull();
  });

  it('should dispatch checked event with value=false when checkbox is checked', async () => {
    let completed = true;
    const handleChecked = vi.fn().mockImplementation((e: CustomEvent) => (completed = e.detail));
    const task: Task = { name: '', description: '', dueDate: '2021-10-10', completed: true };

    const { component } = render(Task, { task });
    const checkbox = screen.getByRole('checkbox');
    component.$on('checked', handleChecked);

    await fireEvent.click(checkbox);

    expect(completed).toEqual(false);
  });

  it('should dispatch checked event with value=true when checkbox is unchecked', async () => {
    let completed = false;
    const handleChecked = vi.fn().mockImplementation((e: CustomEvent) => (completed = e.detail));
    const task: Task = { name: '', description: '', dueDate: '2021-10-10', completed: false };

    const { component } = render(Task, { task });
    const checkbox = screen.getByRole('checkbox');
    component.$on('checked', handleChecked);

    await fireEvent.click(checkbox);

    expect(completed).toEqual(true);
  });
});
