import { ref } from 'vue';

const currentOpen = ref<string | undefined>(undefined);

export const useMoreMenu = () => {
  const open = (id: string) => {
    currentOpen.value = id;
  };

  return {
    lastOpenedId: currentOpen,
    open,
    close,
  };
};
