import { TestBed } from '@angular/core/testing';
import { EventService } from './event.service';

describe('EventService', () => {
  let service: EventService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EventService);
  });

  it('should add an attendee with status when setting status for a non-attendee', (done) => {
    // create a fresh event
    service.create({
      title: 'Test Event A',
      date: '2025-12-01',
      organizerId: 'org-1',
      attendees: [] as any,
    }).subscribe(evt => {
      service.setAttendanceStatus(evt.id, 'guest-1', 'Going').subscribe(ok => {
        expect(ok).toBeTrue();
        service.getInvitedEvents('guest-1').subscribe(list => {
          expect(list.some(e => e.id === evt.id)).toBeTrue();
          const found = list.find(e => e.id === evt.id)!;
          const att = found.attendees.find(a => a.id === 'guest-1');
          expect(att).toBeDefined();
          expect(att?.status).toBe('Going');
          done();
        });
      });
    });
  });

  it('should update status for existing attendee', (done) => {
    service.create({
      title: 'Test Event B',
      date: '2025-12-02',
      organizerId: 'org-2',
      attendees: [{ id: 'peter', status: 'Maybe' } as any],
    }).subscribe(evt => {
      service.setAttendanceStatus(evt.id, 'peter', 'Not Going').subscribe(ok => {
        expect(ok).toBeTrue();
        service.getInvitedEvents('peter').subscribe(list => {
          const found = list.find(e => e.id === evt.id)!;
          const att = found.attendees.find(a => a.id === 'peter');
          expect(att?.status).toBe('Not Going');
          done();
        });
      });
    });
  });
});
