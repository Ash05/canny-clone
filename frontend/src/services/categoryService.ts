import { environment } from '../environments/environment';
import { authService } from './authService';

export interface Category {
  id: number;
  name: string;
}

class CategoryService {
  // Get all categories
  async getAllCategories(): Promise<Category[]> {
    const response = await fetch(`${environment.apiUrl}/categories`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch categories');
    }
    
    return response.json();
  }
}

export const categoryService = new CategoryService();
