import register from "preact-custom-element";
import { useState } from "preact/hooks";
import { signal } from "@preact/signals";

const mutationTracker = signal(0);

const Greeting = ({ count: initialCount }: { count: string }) => {
  const [count, setCount] = useState(Number(initialCount) || 0);
  const increment = () => {
    setCount(count + 1);
    mutationTracker.value += 1;
  };
  // You can also pass a callback to the setter
  const decrement = () => {
    setCount((currentCount) => currentCount - 1);
    mutationTracker.value += 1;
  };

  return (
    <div class=" flex gap-2 bg-pink-200 p-4 my-2 rounded-lg border-2 border-red-300">
      <p class="flex text-lg font-bold text-pink-800 dark:text-pink-800 w-36 items-center">
        Count: {count}
      </p>
      <div class="flex gap-2">
        <button class="button" onClick={increment}>
          +
        </button>
        <button class="button" onClick={decrement}>
          -
        </button>
      </div>
      <div class="flex flex-grow items-center justify-center">
        Total Mutations (Signal): {mutationTracker.value}
      </div>
    </div>
  );
};

register(Greeting, "x-greeting", ["name"], { shadow: false });
