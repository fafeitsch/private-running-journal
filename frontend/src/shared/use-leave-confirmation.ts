import { useConfirm } from "primevue/useconfirm";
import { onBeforeRouteLeave, onBeforeRouteUpdate } from "vue-router";
import {useI18n} from 'vue-i18n';
import {MaybeRef, MaybeRefOrGetter, toValue} from 'vue';

export const useLeaveConfirmation: (dirty: MaybeRefOrGetter<boolean>) => void = (dirty: MaybeRefOrGetter<boolean>) => {
  const confirm = useConfirm();
  const {t} = useI18n()

  onBeforeRouteLeave(() => handleRouteLeave());
  onBeforeRouteUpdate(() => handleRouteLeave());

  function handleRouteLeave(): Promise<boolean> {
    if (!toValue(dirty)) {
      return Promise.resolve(true);
    }
    let resolveFn: (result: boolean) => void;
    const result = new Promise<boolean>((resolve) => (resolveFn = resolve));
    confirm.require({
      group: "leave",
      header: t("shared.confirm.header"),
      accept: () => resolveFn(true),
      reject: () => resolveFn(false),
      message: t("shared.confirm.message"),
      rejectLabel: t("shared.cancel"),
      acceptLabel: t("shared.confirm.discard"),
    });
    return result;
  }
};
