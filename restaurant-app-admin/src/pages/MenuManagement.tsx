import React, { useEffect, useMemo, useState } from 'react';
import { menuApi, type MenuItemDto } from '../services/api';

interface Dish {
  id: string;
  name: string;
  price: number;
  description: string;
  image: string;
  category: string;
  isVip?: boolean;
}

const emptyDishForm = {
  name: '',
  price: '',
  description: '',
  image: '',
  category: '',
};

const getItemId = (item: MenuItemDto) => item.item_id ?? item.itemId ?? '';
const getImageUrl = (item: MenuItemDto) => {
  const image = item.image_url ?? item.imageUrl;

  if (!image) {
    return '/images/default-food.jpg';
  }

  // nếu DB trả về assets/images/...
  if (image.startsWith('assets/')) {
    return image.replace('assets/', '/');
  }

  return image.startsWith('/') ? image : `/${image}`;
};
const mapMenuItem = (item: MenuItemDto): Dish => ({
  id: getItemId(item),
  name: item.name ?? 'Chưa đặt tên',
  price: item.price ?? 0,
  description: item.description ?? '',
  image: getImageUrl(item),
  category: item.category ?? 'Khác',
  isVip: (item.price ?? 0) >= 1000000,
});

const formatVnd = (value: number) => `${Math.round(value).toLocaleString('vi-VN')}đ`;

const MenuManagement: React.FC = () => {
  // Quản lý trạng thái phân loại (Category)
  const [activeCategory, setActiveCategory] = useState<string>('Tất cả');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [dishes, setDishes] = useState<Dish[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  // Quản lý trạng thái Modal
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const [modalMode, setModalMode] = useState<'add' | 'edit'>('add');
  const [editingDish, setEditingDish] = useState<Dish | null>(null);
  const [form, setForm] = useState(emptyDishForm);
  const [isSaving, setIsSaving] = useState<boolean>(false);

  // const categories = useMemo(() => ['Tất cả', ...Array.from(new Set(dishes.map((dish) => dish.category).filter(Boolean)))], [dishes]);
  const [categories, setCategories] = useState<{ id: string; name: string }[]>([]);

  const filteredDishes = useMemo(() => {
    return dishes.filter((dish) => {
      const matchesCategory = activeCategory === 'Tất cả' || dish.category === activeCategory;
      const matchesSearch = dish.name.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesCategory && matchesSearch;
    });
  }, [activeCategory, dishes, searchQuery]);

  const loadMenu = async () => {
    setIsLoading(true);
    setError('');
    try {
      const response = await menuApi.listItems({ page: 1, page_size: 100, keyword: searchQuery });
      setDishes((response.items ?? []).map((item) => mapMenuItem(item, categories)).filter((item) => item.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể tải thực đơn từ máy chủ.');
    } finally {
      setIsLoading(false);
    }
  };

  const loadCategories = async () => {
    try {
      const response = await menuApi.listCategories();
      setCategories(
        response.categories?.map((cat) => ({
          id: cat.category_id ?? '',
          name: cat.name ?? 'Khác',
        })) ?? []
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể tải danh mục.');
    }
  };

  // Thêm categories vào scope hoặc truyền vào
  const mapMenuItem = (item: MenuItemDto, cats: { id: string; name: string }[]): Dish => {
    const catName = cats.find((c) => c.id === (item.category_id ?? item.categoryId))?.name
      ?? item.category
      ?? 'Khác';
    return {
      id: getItemId(item),
      name: item.name ?? 'Chưa đặt tên',
      price: item.price ?? 0,
      description: item.description ?? '',
      image: getImageUrl(item),
      category: catName,
      isVip: (item.price ?? 0) >= 1000000,
    };
  };

  useEffect(() => {
    loadMenu();
    loadCategories();
  }, []);

  const handleOpenModal = (mode: 'add' | 'edit', dish?: Dish) => {
    setModalMode(mode);
    setEditingDish(dish ?? null);
    setForm(
      dish
        ? {
          name: dish.name,
          price: String(dish.price),
          description: dish.description,
          image: dish.image,
          category: dish.category,
        }
        : emptyDishForm,
    );
    setIsModalOpen(true);
  };

  const handleDelete = async (dish: Dish) => {
    if (!window.confirm(`Xóa món "${dish.name}"?`)) return;

    try {
      await menuApi.deleteItem(dish.id);
      setDishes((current) => current.filter((item) => item.id !== dish.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể xóa món ăn.');
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setIsSaving(true);
    setError('');

    const payload = {
      name: form.name.trim(),
      description: form.description.trim(),
      price: Number(form.price) || 0,
      category: form.category.trim(),
      category_id: form.category.trim(),
      image_url: form.image.trim(),
    };

    try {
      if (modalMode === 'edit' && editingDish) {
        const response = await menuApi.updateItem(editingDish.id, payload);
        const updated = response.item ? mapMenuItem(response.item) : { ...editingDish, ...payload, image: payload.image_url };
        setDishes((current) => current.map((dish) => (dish.id === editingDish.id ? updated : dish)));
      } else {
        const response = await menuApi.createItem(payload);
        if (response.item) {
          setDishes((current) => [mapMenuItem(response.item!), ...current]);
        } else {
          await loadMenu();
        }
      }

      setIsModalOpen(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể lưu món ăn.');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="flex flex-col gap-6 animate-fadeIn">
      {/* Header riêng của trang Menu */}
      <div className="flex justify-between items-end border-b border-outline-variant/20 pb-6">
        <div>
          <nav className="flex gap-2 text-on-surface-variant text-[10px] mb-2 uppercase tracking-widest">
            <span>Admin</span>
            <span>/</span>
            <span className="text-primary font-bold">Thực đơn</span>
          </nav>
          <h2 className="font-serif text-5xl font-bold text-on-surface">Quản lý Thực đơn</h2>
        </div>
        <button
          onClick={() => handleOpenModal('add')}
          className="bg-[#735c00] text-white text-xs font-semibold px-6 py-3 rounded-lg shadow-sm hover:shadow-md transition-all active:scale-95 flex items-center gap-2"
        >
          <span className="material-symbols-outlined text-[18px]">add</span>
          Thêm món mới
        </button>
      </div>

      {/* Thanh Tìm kiếm & Phân loại bộ lọc */}
      <div className="flex flex-wrap justify-between items-center gap-6 my-4">
        {/* Category Tabs */}
        <div className="relative flex gap-8 border-b border-outline-variant/30 pb-px">
          {['Tất cả', ...categories.map((c) => c.name)].map((categoryName) => (
            <button
              key={categoryName}
              onClick={() => setActiveCategory(categoryName)}
              className={`text-xs font-semibold py-2 relative transition-colors ${activeCategory === categoryName ? 'text-[#735c00]' : 'text-on-surface-variant hover:text-[#735c00]'
                }`}
            >
              {categoryName}
              {activeCategory === categoryName && (
                <span className="absolute bottom-0 left-0 w-full h-[2px] bg-[#735c00]"></span>
              )}
            </button>
          ))}
        </div>

        {/* Ô tìm kiếm */}
        <div className="relative">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px]">search</span>
          <input
            className="pl-10 pr-4 py-2 bg-[#f3f4f5] border border-outline-variant rounded-full text-sm w-64 focus:outline-none focus:border-[#735c00] transition-all"
            placeholder="Tìm kiếm món ăn..."
            type="text"
            value={searchQuery}
            onChange={(event) => setSearchQuery(event.target.value)}
          />
        </div>
      </div>

      {error && <p className="text-sm text-red-600">{error}</p>}
      {isLoading && <p className="text-sm text-on-surface-variant">Đang tải thực đơn từ API...</p>}

      {/* Lưới danh sách món ăn */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {filteredDishes.map((dish) => (
          <div key={dish.id} className="group relative bg-white rounded-xl shadow-sm border border-outline-variant/10 overflow-hidden hover:shadow-lg transition-all duration-300 flex flex-col">
            <div className="relative h-48 overflow-hidden bg-gray-100">
              {dish.isVip && (
                <div className="absolute top-0 left-0 bg-[#735c00] px-3 py-1 text-white text-[10px] z-10 font-bold uppercase tracking-widest rounded-br-lg">
                  VIP Selection
                </div>
              )}
              <img alt={dish.name} className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110" src={dish.image} />

              {/* Lớp phủ chứa nút Sửa/Xóa khi hover */}
              <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center gap-4 transition-opacity duration-300">
                <button
                  onClick={() => handleOpenModal('edit', dish)}
                  className="w-12 h-12 bg-white text-on-surface rounded-full flex items-center justify-center hover:bg-[#735c00] hover:text-white transition-all shadow-lg"
                >
                  <span className="material-symbols-outlined">edit</span>
                </button>
                <button onClick={() => handleDelete(dish)} className="w-12 h-12 bg-white text-red-600 rounded-full flex items-center justify-center hover:bg-red-600 hover:text-white transition-all shadow-lg">
                  <span className="material-symbols-outlined">delete</span>
                </button>
              </div>
            </div>

            <div className={`p-6 flex-grow ${dish.isVip ? 'border-t-4 border-[#735c00]' : ''}`}>
              <div className="flex justify-between items-start mb-2">
                <h3 className="font-serif text-xl font-bold text-on-surface">{dish.name}</h3>
                <span className="text-sm font-bold text-[#735c00]">{formatVnd(dish.price)}</span>
              </div>
              <p className="text-on-surface-variant text-sm line-clamp-2">{dish.description}</p>
            </div>
          </div>
        ))}

        {/* Nút Khung giữ chỗ Thêm Món Nhanh */}
        <div
          onClick={() => handleOpenModal('add')}
          className="border-2 border-dashed border-outline-variant/50 rounded-xl flex flex-col items-center justify-center p-12 group cursor-pointer hover:border-[#735c00]/50 hover:bg-[#edeeef] transition-all"
        >
          <div className="w-16 h-16 rounded-full bg-[#e1e3e4] text-on-surface-variant flex items-center justify-center group-hover:bg-[#735c00] group-hover:text-white transition-all mb-4">
            <span className="material-symbols-outlined text-[32px]">add_circle</span>
          </div>
          <p className="text-xs font-semibold text-on-surface-variant group-hover:text-[#735c00] transition-colors">Thêm món ăn mới vào thực đơn</p>
        </div>
      </div>

      {/* Unified Modal cho việc Add/Edit món ăn */}
      {isModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm transition-opacity duration-300">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-2xl overflow-hidden transform scale-100 transition-transform duration-300">
            <div className="px-8 py-6 border-b border-outline-variant/20 flex justify-between items-center">
              <h3 className="font-serif text-xl font-bold text-on-surface">
                {modalMode === 'add' ? 'Thêm món mới' : 'Cập nhật món ăn'}
              </h3>
              <button className="text-on-surface-variant hover:text-[#735c00] transition-colors" onClick={() => setIsModalOpen(false)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>

            <form className="p-8 space-y-6" onSubmit={handleSubmit}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Tên món ăn</label>
                  <input className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface" placeholder="Ví dụ: Bò Wagyu A5..." type="text" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} required />
                </div>
                <div className="space-y-2">
                  <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Giá (VNĐ)</label>
                  <input className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface" placeholder="0" type="number" value={form.price} onChange={(event) => setForm({ ...form, price: event.target.value })} required />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Danh mục</label>
                <input className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface" placeholder="Ví dụ: Món chính" type="text" value={form.category} onChange={(event) => setForm({ ...form, category: event.target.value })} required />
              </div>
              <div className="space-y-2">
                <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Mô tả chi tiết</label>
                <textarea className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface resize-none" placeholder="Mô tả thành phần, cách chế biến..." rows={3} value={form.description} onChange={(event) => setForm({ ...form, description: event.target.value })}></textarea>
              </div>
              <div className="space-y-2">
                <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">URL hình ảnh món ăn</label>
                <input className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface" placeholder="https://..." type="url" value={form.image} onChange={(event) => setForm({ ...form, image: event.target.value })} />
              </div>
              <div className="flex justify-end items-center gap-6 pt-4">
                <button className="text-xs font-semibold text-on-surface-variant hover:text-[#735c00] transition-colors px-4 py-2" onClick={() => setIsModalOpen(false)} type="button">Hủy bỏ</button>
                <button disabled={isSaving} className="px-10 py-3 bg-[#735c00] text-white text-xs font-bold rounded-lg shadow-md hover:bg-[#735c00]/90 transition-all active:scale-95 uppercase tracking-widest disabled:opacity-60" type="submit">
                  {isSaving ? 'Đang lưu...' : modalMode === 'add' ? 'Tạo món mới' : 'Lưu thay đổi'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default MenuManagement;
