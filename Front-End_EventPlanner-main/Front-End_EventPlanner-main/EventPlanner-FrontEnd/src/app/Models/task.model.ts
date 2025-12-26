export interface TaskModel {
  id: string;
  title: string;
  description?: string;
  date?: string; // ISO date
  createdBy?: string; // user id
  assigneeId?: string; // user id
}
