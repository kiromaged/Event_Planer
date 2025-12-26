import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class UsersService {
  // Simple static user directory for demo purposes. Replace with API-backed user lookup later.
  private users: Record<string, { id: string; name: string }> = {
    'user-1': { id: 'user-1', name: 'Sami Elbialley' },
    'user-2': { id: 'user-2', name: 'Alice Johnson' },
    'user-3': { id: 'user-3', name: 'Bob Smith' }
  };

  getDisplayName(userId: string | null | undefined): string {
    if (!userId) return 'Unknown';
    return this.users[userId]?.name ?? userId;
  }
}
