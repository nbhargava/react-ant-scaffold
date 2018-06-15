import { EventSubscription } from 'fbemitter';

export function removeListener(listener: EventSubscription | null) {
  if (listener !== null) {
    listener.remove();
  }
}

