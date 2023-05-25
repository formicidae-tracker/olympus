import { Component, inject } from '@angular/core';
import { MatSnackBarRef } from '@angular/material/snack-bar';

@Component({
  selector: 'app-snack-network-offline',
  templateUrl: './snack-network-offline.component.html',
  styleUrls: ['./snack-network-offline.component.scss'],
})
export class SnackNetworkOfflineComponent {
  snackBarRef = inject(MatSnackBarRef);
}
