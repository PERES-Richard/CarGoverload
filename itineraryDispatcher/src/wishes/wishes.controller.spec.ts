import { Test, TestingModule } from '@nestjs/testing';
import { WishesController } from './wishes.controller';

describe('WishesController', () => {
  let controller: WishesController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [WishesController],
    }).compile();

    controller = module.get<WishesController>(WishesController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
