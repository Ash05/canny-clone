import React from 'react';
import { useForm } from 'react-hook-form';

function Boards() {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm();

  const boardName = watch('boardName');

  const validateBoardName = async (name) => {
    if (/[^a-zA-Z0-9 ]/.test(name)) {
      return 'Board name cannot contain special characters.';
    }
    try {
      const res = await fetch('http://localhost:8080/validate-board', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
      });
      const data = await res.json();
      if (!data.isUnique) {
        return 'Board name must be unique.';
      }
    } catch (error) {
      return 'Error validating board name.';
    }
    return true;
  };

  const onSubmit = async (data) => {
    try {
      const res = await fetch('http://localhost:8080/create-board', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: data.boardName }),
      });
      if (res.ok) {
        alert('Board created successfully!');
      } else {
        alert('Error creating board.');
      }
    } catch (error) {
      alert('Error creating board.');
    }
  };

  return (
    <div>
      <h1>Boards</h1>
      <form onSubmit={handleSubmit(onSubmit)}>
        <input
          type="text"
          {...register('boardName', {
            required: 'Board name cannot be empty.',
            validate: validateBoardName,
          })}
          placeholder="Enter board name"
        />
        {errors.boardName && <p style={{ color: 'red' }}>{errors.boardName.message}</p>}
        <button type="submit">Create Board</button>
      </form>
    </div>
  );
}

export default Boards;
